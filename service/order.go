package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"myproject/dao"
	"myproject/model"
	"strconv"
)

func CreateOrder(userID uint, products []model.Product) (string, error) {
	db := dao.DB

	//查询用户是否存在,并将用户订单信息放到order中
	// 查询用户是否存在
	var user model.User
	result := db.Where("id = ?", userID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("user with ID %d not found", userID)
	}
	if result.Error != nil {
		return "", result.Error
	}

	order := model.Order{
		UserID:     userID,
		TotalPrice: 0.0,
		User:       user,
	}
	//orders := model.Order{
	//	UserId:     userID,
	//	TotalPrice: 0.0,
	//}
	//遍历商品计算总价格
	for _, product := range products {
		order.TotalPrice += product.Price
	}
	//在数据库中创建订单
	result = db.Create(&order)
	if result.Error != nil {
		return strconv.Itoa(int(order.ID)), result.Error
	}

	for _, product := range products {
		db.Model(&order).Create(&product)
	}

	return strconv.Itoa(int(order.ID)), result.Error
}
