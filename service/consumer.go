package service

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"log"
	"myproject/dao"
	"myproject/model"
)

// StartStockConsumer 启动kafka消息队列生产者
func StartStockConsumer() {
	config := sarama.NewConfig()         //创建 Kafka 消费者配置对象
	config.Consumer.Return.Errors = true //配置消费者返回错误信息，方便调试和处理异常

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config) //创建一个 Kafka 消费者，连接到 localhost:9092
	if err != nil {
		log.Fatal("Failed to create consumer:", err) //如果初始化失败，通过 log.Fatal 抛出异常并终止程序
	}
	defer func(consumer sarama.Consumer) {
		err := consumer.Close() //确保在函数退出时关闭 Kafka 消费者连接，释放资源
		if err != nil {

		}
	}(consumer)

	//消费指定主题（stock_updates）的特定分区（分区编号为 0）的消息
	partitionConsumer, err := consumer.ConsumePartition("stock_updates", 0, sarama.OffsetNewest)
	//sarama.OffsetNewest: 从最新的消息开始消费。sarama.OffsetOldest从最早的消息开始消费
	if err != nil {
		log.Fatal("Failed to consume partition:", err)
	}
	defer func(partitionConsumer sarama.PartitionConsumer) {
		err := partitionConsumer.Close() //确保在函数退出时关闭分区消费者连接，释放资源
		if err != nil {

		}
	}(partitionConsumer)

	for {
		select {
		case msg := <-partitionConsumer.Messages(): //从 Kafka 分区中读取消息
			var stockMsg StockMessage
			if json.Unmarshal(msg.Value, &stockMsg) == nil { //将 Kafka 消息的内容（msg.Value）反序列化为 StockMessage 结构体
				updateMySQLStock(stockMsg) //调用 updateMySQLStock 函数，根据消息内容更新 MySQL 数据库中的库存
			}
		case err := <-partitionConsumer.Errors(): //从 Kafka 分区中读取错误信息
			log.Printf("Consumer error: %v", err) //使用Printf记录错误信息
		}
	}
}

// 更新MySQL库存
func updateMySQLStock(msg StockMessage) {
	dao.DB.Model(&model.Product{}).Where(
		"id = ?", msg.ProductID).Update(
		"number", gorm.Expr("number - ?", msg.Quantity)) //计算新的库存值
}
