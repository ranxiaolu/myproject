package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	gorm.Model
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Avatar       string    `json:"avatar"`
	Nickname     string    `json:"nickname"`
	Introduction string    `json:"introduction"`
	Telephone    string    `json:"telephone"`
	QQ           string    `json:"qq"`
	Gender       string    `json:"gender"`
	Email        string    `json:"email"`
	Birthday     time.Time `json:"birthday"`
}

//m := map[string]interface{
//	"nickname":User.Nickname,
//	"telephone":User.Telephone,
//	"qq":User.QQ,
//	"gender":User.Gender,
//	"introduction":User.Introduction,
//	"email":User.Email,
//	"birthday":User.Birthday,
//}
