package dao

import (
	"errors"
	"gorm.io/gorm"
	"myproject/model"
)

//func GetDB(*gorm.DB) *gorm.DB {
//	return db
//}

// FindUser 用户注册
//

func FindUser(user *model.User) (bool, bool) {
	result := DB.Model(user).Where("username = ?", user.Username).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, true //用户不存在，查询成功
	} else if result.Error != nil {
		return false, false
	} //查询失败
	return true, true //用户存在，查询操作成功
}

//func FindUser(user *model.User) error {
//	if user == nil {
//		return fmt.Errorf("user 指针为 nil")
//
//	}
//	result := DB.Where("username=?", user.Username).First(&user).Error
//	if result == nil {
//		return nil
//	}
//	return result
//}

// CreateUser 创建新账号
func CreateUser(username, password string) error {
	var user model.User
	user.Username = username
	user.Password = password
	//传入数据库
	result := DB.Create(&user)
	return result.Error
}
