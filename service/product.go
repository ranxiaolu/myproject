package service

import (
	"errors"
	"myproject/dao"
	"myproject/model"
)

func GetProducts() ([]model.Product, error) {
	db := dao.DB
	var products []model.Product
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
func GetProductDetails(Name uint) (model.Product, error) {
	db := dao.DB
	var product model.Product
	result := db.First(&product, "name = ?", Name)
	if result.Error != nil {
		return product, result.Error
	}
	return product, nil
}
func AddProductTOCart(productID uint, quantity int) error {
	var cartItem model.CartItem
	db := dao.DB
	var product model.Product
	if result := db.First(&product, "ID=?", productID); result.Error != nil {
		return errors.New("product not found")
	}
	//err := dao.AddProductTOCart(productID).Error
	//if err != nil {
	//	return errors.New("product not found")
	//}
	cartItem.Quantity += quantity
	result := db.Save(&cartItem)
	if result.Error != nil {
		return result.Error
	}
	return nil

}
func GetCartList(userID string) ([]model.CartItem, error) {
	Cartlist, err := dao.GetCartList(userID)
	//出错
	if err != nil {
		return nil, err
	}
	return Cartlist, nil
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
