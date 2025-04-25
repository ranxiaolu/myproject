package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"myproject/dao"
	"myproject/model"
	"time"
)

// RegisterUser 用户注册
func RegisterUser(user model.User) error {
	db := dao.DB

	//传入数据库
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func CreateToken(username string, password string) (string, string, error) {
	// 生成JWT令牌
	claims := &model.CustomClaims{
		Username: username,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 令牌有效期1小时
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(model.JwtSecretKey)) //嵌套结构不符合标准
	if err != nil {                                                    //jwt.SigningMethodHS256要求密钥为[]byte类型
		return "", "", err
	}

	// 生成刷新令牌
	refreshClaims := &model.CustomClaims{
		Username: username,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 刷新令牌有效期1天
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(model.JwtSecretKey))
	if err != nil {
		return "", "", err
	}

	return refreshTokenString, tokenString, nil
}

// LoginUser 用户登录
func LoginUser(username, password string) (string, error) {
	db := dao.DB
	var user model.User
	//用户是否存在
	err := db.Where("username=?", username).First(&user).Error
	if err != nil {
		return "", errors.New("wrong ")
	}
	//验证密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("password wrong ")
	}
	return "", nil

}

// RefreshToken 刷新token
func RefreshToken(tokenString string) (string, error) {
	//解析JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //ok=false 不是预期的签名方法
			return nil, errors.New("unexpected signing method")
		}

		return tokenString, nil
	})
	//token解析出错
	if err != nil {
		return "", err
	}
	//验证JWT
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//是否过期
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return "", errors.New("token has expired")
		}

		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": claims["username"],
			"exp":      time.Now().Add(time.Hour * 24).Unix(), // 新 Token 的有效期为 24 小时
		})

		newTokenString, err := newToken.SignedString("new-jwt-secret")
		if err != nil {
			return "", err
		}
		return newTokenString, nil
	}

	return "", errors.New("refresh wrong")
}

// ChangePassword 用户修改密码
func ChangePassword(username string, oldPassword string, newPassword string) error {
	db := dao.DB
	//查询用户是否存在
	var user model.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	//// 比较用户密码和用户输入的旧密码是否一致   数据库密码并未加密，导致冲突
	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
	//	return errors.New("password wrong")
	//}
	//// 更新密码 加密新密码
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	//if err != nil {
	//	return err
	//}
	//user.Password = string(hashedPassword)
	////保存到数据库
	//err = db.Model(&user).Where("username=?", user.Username).Update("password", newPassword).Error
	//if err != nil {
	//	return fmt.Errorf("数据库保存失败: %w", err)
	//}
	// 比较用户密码和用户输入的旧密码是否一致

	//// 更新密码 加密新密码
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	//if err != nil {
	//	return err
	//}
	//user.Password = string(hashedPassword)

	if user.Password != oldPassword {
		return errors.New("旧密码错误")
	}

	// 保存到数据库
	err := db.Model(&user).Where("username=?", user.Username).Update("password", newPassword).Error
	if err != nil {
		return fmt.Errorf("数据库保存失败: %w", err)
	}

	return nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(userID string) (model.User, error) {
	db := dao.DB
	var user model.User
	result := db.Where("ID = ?", userID).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// UpdateUserInfo 修改用户信息
func UpdateUserInfo(user *model.User) error {
	db := dao.DB

	result := db.Model(user).Where(
		"id = ? AND username = ? AND password=?", user.ID, user.Username, user.Password).Omit(
		"id,password,username").Updates(user)

	// 1. 处理数据库错误
	if result.Error != nil {
		return fmt.Errorf("数据库操作失败: %w", result.Error)
	}

	// 2. 处理无行更新的情况
	if result.RowsAffected == 0 {
		// 检查是否存在符合条件的用户（确认是密码错误还是数据未变更）
		var exists bool
		if err := db.Model(&model.User{}).
			Select("count(*) > 0").
			Where("user_id = ? AND password = ?", user.ID, user.Password).
			Find(&exists).Error; err != nil {

			return fmt.Errorf("failed to verify user existence: %w", err)
		}

		if !exists {
			return fmt.Errorf("invalid credentials: user ID or username or password incorrect")
		} else {
			// 用户存在但数据未变更
			return fmt.Errorf("no changes detected")
		}
	}

	return nil
}
