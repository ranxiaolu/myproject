package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"myproject/dao"
	"myproject/model"
)

// 创建 Kafka 生产者配置对象
var producer sarama.SyncProducer

func InitKafka() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //要求 Kafka 在所有副本都确认消息写入后才返回成功响应
	config.Producer.Retry.Max = 5                    //允许最多重试 5 次
	config.Producer.Return.Successes = true          //返回成功发送的消息信息

	var err error
	producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, config) //创建一个同步 Kafka 生产者，连接到 localhost:9092
	if err != nil {
		log.Fatal("Failed to init kafka producer:", err) //如果初始化失败，程序会通过 log.Fatal 抛出异常并终止运行
	}
}

type StockMessage struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func CreateOrder(userID uint, items []model.OrderItem) (uint, error) {
	ctx := context.Background()
	db := dao.DB

	//  检查用户是否存在
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return 0, fmt.Errorf("用户不存在")
	}

	//  预扣Redis库存
	for _, item := range items {
		stockKey := fmt.Sprintf("stock:%d", item.ProductID)
		currentStock, err := dao.Redis.DecrBy(ctx, stockKey, int64(item.Quantity)).Result()
		if err != nil {
			return 0, fmt.Errorf("库存操作失败")
		}
		//如果 Redis 中的库存不足（currentStock < 0），则回滚库存操作（通过 IncrBy 增加库存），并返回错误
		if currentStock < 0 {
			// 库存不足，回滚并返回错误
			dao.Redis.IncrBy(ctx, stockKey, int64(item.Quantity))
			return 0, fmt.Errorf("商品%d库存不足", item.ProductID)
		}
	}

	// 创建数据库事务
	//使用 GORM 的事务功能，开始一个数据库事务
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			//如果在事务中发生 panic（例如数据库操作失败），则回滚事务，并调用 rollbackRedisStock 函数回滚 Redis 中的库存。
			tx.Rollback()
			// 回滚Redis库存
			rollbackRedisStock(items)
		}
	}()

	// 创建订单
	order := model.Order{UserID: userID}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		rollbackRedisStock(items)
		return 0, err
	}

	// 创建订单项并计算总价
	totalPrice := 0.0
	for _, item := range items {
		product, err := GetProductDetails(item.ProductID) // 带缓存的查询
		if err != nil {
			//如果创建失败，回滚事务并回滚 Redis 库存
			tx.Rollback()
			rollbackRedisStock(items)
			return 0, err
		}

		orderItem := model.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			//Price:     product.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			//如果任何步骤失败，回滚事务并回滚 Redis 库存
			tx.Rollback()
			rollbackRedisStock(items)
			return 0, err
		}

		totalPrice += product.Price * float64(item.Quantity)
	}

	// 更新订单总价
	if err := tx.Model(&order).Update("total_price", totalPrice).Error; err != nil {
		tx.Rollback()
		rollbackRedisStock(items)
		return 0, err
	}

	//  提交事务
	if err := tx.Commit().Error; err != nil {
		//如果提交失败，回滚 Redis 库存
		rollbackRedisStock(items)
		return 0, err
	}

	// 发送Kafka消息异步更新MySQL
	go asyncUpdateStock(items) //在一个独立的 goroutine 中调用 asyncUpdateStock 函数，将库存更新的消息发送到 Kafka

	return order.ID, nil
}

// 用于在事务失败时回滚库存操作
func rollbackRedisStock(items []model.OrderItem) {
	ctx := context.Background()
	for _, item := range items {
		dao.Redis.IncrBy(ctx,
			fmt.Sprintf("stock:%d", item.ProductID),
			int64(item.Quantity),
		)
	}
}

// 异步更新库存
func asyncUpdateStock(items []model.OrderItem) {
	for _, item := range items {
		//构造 Kafka 消息，包含商品 ID 和数量
		msg := StockMessage{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		jsonMsg, _ := json.Marshal(msg)

		//将消息发送到 Kafka 的 stock_updates 主题
		_, _, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: "stock_updates",
			Value: sarama.ByteEncoder(jsonMsg),
		})
		//如果发送失败，记录日志
		if err != nil {
			log.Printf("Failed to send stock update: %v", err)

		}
	}
}

//func CreateOrder(userID uint, items []model.OrderItem) (uint, error) {
//	ctx := context.Background()
//	db := dao.DB
//
//	//查询用户是否存在,并将用户订单信息放到order中
//	// 查询用户是否存在
//	var user model.User
//	result := db.Where("id = ?", userID).First(&user)
//	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//		return 400, fmt.Errorf("user with ID %d not found", userID)
//	}
//	if result.Error != nil {
//		return 400, result.Error
//	}
//
//	order := model.Order{
//		UserID:     userID,
//		TotalPrice: 0.0,
//		User:       user,
//		//OrderItems: ,
//	}
//
//	//遍历商品计算总价格
//	for _, item := range items {
//		product, err := GetProductDetails(item.ID)
//		if err != nil {
//			return 400, err
//		}
//		// 购物列表总价格
//		order.TotalPrice += product.Price * float64(item.Quantity)
//		product.Number = product.Number - item.Quantity
//		// 更新库存
//		db.Model(&product).Updates("number")
//	}
//	//在数据库中创建订单
//	result = db.Create(&order)
//	if result.Error != nil {
//		return 400, result.Error
//	}
//
//	//for _, product := range products {
//	//	db.Model(&order.OrderItems).Updates(&product)
//	//}
//
//	return order.ID, nil
//}
