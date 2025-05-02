package model

import "gorm.io/gorm"

// import (
//
//	"gorm.io/gorm"
//
// )
//
// // CartItem 购物车条目模型
//
//	type CartItem struct {
//		gorm.Model
//		UserID    uint64 `gorm:"notNull" json:"user_id"`
//		ProductID uint64 `gorm:"notNull" json:"product_id"`
//		Quantity  uint   `gorm:"notNull" json:"quantity"`
//	}
//
// // TableName 返回数据库表名
//
//	func (CartItem) TableName() string {
//		return "cart_items"
//	}
type CartItem struct {
	gorm.Model         //`json:"gorm_._model"`
	UserID     uint    `gorm:"index;notNull" json:"user_id" `   //关联用户
	ProductID  uint    `gorm:"index;notNull" json:"product_id"` //关联产品
	Quantity   int     `json:"quantity" json:"quantity,omitempty"`
	User       User    `gorm:"foreignKey:UserID" `
	Product    Product `gorm:"foreignKey:ProductID"` //关联产品模型
}
