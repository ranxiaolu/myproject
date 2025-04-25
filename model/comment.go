package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ProductID    uint    `gorm:"index;notNull;" json:"product_id"`
	UserID       uint    `gorm:"index;notNull" json:"user_id"`
	Content      string  `json:"content"`
	CommentModel int     `json:"comment_model"` //(1 点赞 2 点踩）
	PraiseCount  int     `json:"praise_count"`
	BadModel     int     `json:"dislike_count"`
	User         User    `gorm:"foreignKey:UserID"`
	Product      Product `gorm:"foreignKey:ProductID"` // 关联产品模型
}
