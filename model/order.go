package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint        `gorm:"index;notNull;foreignKey" json:"user_id"`
	Address    string      `json:"address"`
	TotalPrice float64     `json:"total"`
	User       User        `gorm:"foreignKey:UserID"` // 关联用户模型
	OrderItems []OrderItem `json:"-" gorm:"OrderID"`
}
type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"index;noyNull;" json:"order_id"`
	ProductID uint    `gorm:"index;notNull" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID"` // 关联产品模型
	Quantity  int     `json:"quantity"`
}
