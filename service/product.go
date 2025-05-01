package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"myproject/dao"
	"myproject/model"
	"time"
)

func GetProducts() ([]model.Product, error) {
	db := dao.DB
	var products []model.Product //商品列表
	result := db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func SearchProduct(productName string) (model.Product, error) {
	db := dao.DB
	var product model.Product
	result := db.Where("name = ?", productName).Find(&product)
	if result.Error != nil {
		return product, result.Error
	}
	return product, nil
}
func GetProductDetails(productID uint) (model.Product, error) {
	/*
		查询 Redis 缓存:
		如果缓存命中（cacheData 不为空且反序列化成功），直接返回商品信息。
		缓存未命中:
		查询 MySQL 数据库，获取商品信息。
		如果数据库查询失败，返回错误。
		写入 Redis 缓存:
		将查询到的商品信息序列化为 JSON 数据，并写入 Redis 缓存。
		返回结果:
		返回查询到的商品信息。
	*/
	ctx := context.Background()                      //创建上下文对象
	cacheKey := fmt.Sprintf("product:%d", productID) //构造 Redis 缓存的键，格式为 "product:商品ID"

	// 先查Redis缓存
	cacheData, err := dao.Redis.Get(ctx, cacheKey).Result() //从 Redis 中获取缓存数据，cacheKey包含商品的 JSON 数据
	if err == nil {
		var product model.Product
		if json.Unmarshal([]byte(cacheData), &product) == nil { //将 JSON 数据反序列化为 model.Product 结构体
			//如果反序列化成功，直接返回商品信息，避免查询数据库
			return product, nil
		}
	}

	// 如果 Redis 缓存中没有找到商品信息，则从 MySQL 数据库中查询。
	var product model.Product
	if err := dao.DB.First(&product, productID).Error; err != nil {
		return product, err
	}

	// 写入Redis缓存
	productJSON, _ := json.Marshal(product)                  //将商品信息序列化为 JSON 数据
	dao.Redis.Set(ctx, cacheKey, productJSON, 5*time.Minute) //将商品信息写入 Redis 缓存
	//5*time.Minute 设置缓存的过期时间为 5 分钟，以避免缓存数据过时

	return product, nil
}

//func GetProductDetails(ProductID uint) (model.Product, error) {
//	db := dao.DB
//	var product model.Product
//	result := db.Where("id = ?", ProductID).First(&product)
//	if result.Error != nil {
//		return product, result.Error
//	}
//	return product, nil
//}

func AddProductTOCart(productID uint, quantity int) error {
	var cartItem model.CartItem
	db := dao.DB
	var product model.Product
	if result := db.Where("id = ?", productID).First(&product); result.Error != nil {
		return errors.New("product not found")
	}

	cartItem.Quantity += quantity
	result := db.Save(&cartItem)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func GetCartList(userID string) ([]model.CartItem, error) {
	cartItem, err := dao.GetCartList(userID)
	//出错
	if err != nil {
		return nil, err
	}
	return cartItem, nil
}

func GetProductsByType(productType string) ([]model.Product, error) {
	var products []model.Product
	db := dao.DB
	result := db.Where("type = ?", productType).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
