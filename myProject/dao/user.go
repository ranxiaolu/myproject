package dao

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"myproject/model"
)

var db *gorm.DB

func init() {
	dsn := "root:5201314@tcp(127.0.0.1:3306)/myshop?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false, // 启用外键约束
	})
	if err != nil {
		panic("failed to connect database")
	}
	//自动迁移数据库
	// 迁移时需按依赖顺序执行
	err = db.AutoMigrate(
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
	//检查 GORM 日志以获取详细错误信息：
	db.Logger = db.Logger.LogMode(logger.Info)

}

func GetDB() *gorm.DB {
	return db
}

// RegisterUser 用户注册
func RegisterUser(user model.User) error {
	//查询用户是否存在
	db := GetDB()
	var registeredUser model.User
	err := db.Where("username=?", user.Username).First(&registeredUser).Error
	if err == nil {
		return errors.New("user already exists")
	}
	return nil
}

// CreateUser 创建新账号
func CreateUser(username, password string) error {
	var user model.User
	user.Username = username
	user.Password = password
	//传入数据库
	result := db.Create(&user)
	return result.Error
}
