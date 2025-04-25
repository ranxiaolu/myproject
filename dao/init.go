package dao

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"myproject/model"
)

var DB *gorm.DB
var ctx = context.Background()

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
	// 连接 Redis
	rdb := redis.NewClient(&redis.Options{ //创建 Redis 客户端
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// 确保连接成功
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MySQL and Redis connected successfully!")
}
