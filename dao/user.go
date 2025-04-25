package dao

import (
	"errors"
	"gorm.io/gorm"
	"myproject/model"
)

func FindUser(user *model.User) (bool, bool) {
	result := DB.Model(user).Where("username = ?", user.Username).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, true //用户不存在，查询成功
	} else if result.Error != nil {
		return false, false
	} //查询失败
	return true, true //用户存在，查询操作成功
}

// CreateUser 创建新账号
func CreateUser(username, password string) error {
	var user model.User
	user.Username = username
	user.Password = password
	//传入数据库
	result := DB.Create(&user)
	return result.Error
}
