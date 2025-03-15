package model

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID    uint    `gorm:"index;notNull" json:"user_id" `   //关联用户
	ProductID uint    `gorm:"index;notNull" json:"product_id"` //关联产品
	Quantity  int     `json:"quantity" json:"quantity,omitempty"`
	User      User    `gorm:"foreignKey:UserID" `
	Product   Product `gorm:"foreignKey:ProductID"` //关联产品模型
}
