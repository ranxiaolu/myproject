package dao

import "myproject/model"

func GetCartList(userID string) ([]model.CartItem, error) {
	db := DB
	var cartItems []model.CartItem

	//将数据库中的查询结果放到cartItems中
	result := db.Where("user_id = ?", userID).Find(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return cartItems, nil
}
