package model

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username     string `json:"username"`
	Password     string `json:"password"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Telephone    string `json:"telephone"`
	QQ           string `json:"qq"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
}
