package dao

import "myproject/model"

//	func GetProducts() ([]model.Product, error) {
//		db := GetDB()
//		var products []model.Product
//		result := db.Find(&products)
//		if result.Error != nil {
//			return nil, result.Error
//		}
//		return products, nil
//	}
//
// 查看商品是否存在
func AddProductTOCart(productID string) error {
	db := GetDB()
	var product model.Product
	//查询商品是否存在
	result := db.First(&product, "product_id = ?", productID)

	if result.RowsAffected == 0 {

		return db.Create(&product).Error
	}
	return nil
}
func GetCartList(userID string) ([]model.CartItem, error) {
	db := GetDB()
	var cartItems []model.CartItem

	//将数据库中的查询结果放到cartItems中
	result := db.Where("user_id = ?", userID).Find(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return cartItems, nil
}
