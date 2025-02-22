package model

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Price       float64   `json:"price"`
	CommentNum  int       `json:"comment_num"`
	IsAddedCart string    `json:"is_addedCart"`
	Cover       string    `json:"cover"`
	PublishTime time.Time `json:"publish_time"`
	Link        string    `json:"link"`
}
