package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"myproject/model"
)

var DB *gorm.DB
var ctx = context.Background()
var Redis *redis.Client // 添加全局Redis客户端，存储redis客户端实例

func Init() {
	dsn := "root:5201314@tcp(127.0.0.1:3306)/myshop?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info), // 启用详细日志
		DisableForeignKeyConstraintWhenMigrating: false,                               // 启用外键约束
	})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	//自动迁移数据库
	// 迁移时需按依赖顺序执行
	err = DB.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.CartItem{},
		&model.Comment{},
	)
	if err != nil {
		panic(err) // 避免静默失败
	}
	if DB == nil {
		panic("database not initialized") // 确保全局变量已初始化
	}
	//检查 GORM 日志以获取详细错误信息：
	DB.Logger = DB.Logger.LogMode(logger.Info)

	//// 连接 Redis
	//rdb := redis.NewClient(&redis.Options{ //创建 Redis 客户端
	//	Addr:     "localhost:6379",
	//	Password: "",
	//	DB:       0, // 使用默认的数据库
	//})

	//// Ping测试确保连接成功
	//_, err = rdb.Ping(ctx).Result()
	//if err != nil {
	//	log.Fatal(err)
	//}
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 如果 Redis 有密码需填写
		DB:       0,
	})

	// 检查 Redis 连接
	if _, err := Redis.Ping(ctx).Result(); err != nil {
		panic("failed to connect redis: " + err.Error())
	}

	// 初始化库存缓存（服务启动时执行）
	initStockCache()
	fmt.Println("MySQL and Redis connected successfully!")

}

// initStockCache 将 MySQL 中的商品库存数据同步到 Redis 缓存中
func initStockCache() {
	var products []model.Product
	if err := DB.Find(&products).Error; err == nil {

		for _, p := range products {
			// 使用set方法设置键及相应值
			Redis.Set(ctx, fmt.Sprintf("stock:%d", p.ID), p.Number, 0) // 永久存储
		}
	}
}
