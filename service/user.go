package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"myproject/dao"
	"myproject/model"
	"time"
)

// RegisterUser 用户注册
func RegisterUser(user model.User, password string) error {
	if err := dao.RegisterUser(user); err != nil {
		return errors.New("user already exists")
	}
	return dao.CreateUser(user.Username, password)
}

// LoginUser 用户登录
func LoginUser(username, password string) (string, error) {
	db := dao.GetDB()
	var user model.User
	err := db.Where("username=?", username).First(&user).Error
	if err != nil {
		return "", errors.New("wrong ")
	}
	//验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong ")
	}
	//生成JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), //有效期24小时
	})

	tokenString, err := token.SignedString("your-jwt-secret")
	if err != nil {
		return "", err
	}
	return tokenString, nil
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
func ChangePassword(user model.User, oldPassword string, newPassword string) error {
	db := dao.GetDB()
	//查询用户是否存在
	var count int64
	err := db.Model(&model.User{}).Where("username=?", user.Username).Count(&count).Error
	//查询过程是否出错
	if err != nil {
		return err
	}
	//用户是否存在
	if count == 0 {
		return errors.New("user not found")
	}

	// 比较旧密码和用户输入的旧密码是否一致
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("password wrong")
	}

	// 更新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	//保存到数据库
	err = db.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(userID string) (model.User, error) {
	db := dao.GetDB()
	var user model.User
	result := db.First(&user, "user_id = ?", userID)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// UpdateUserInfo 修改用户信息
func UpdateUserInfo(user model.User, info map[string]interface{}) error {
	db := dao.GetDB()
	err := db.First(&user, "user_id = ?", user.ID)
	if err != nil {
		return errors.New("username not found")
	}
	//修改信息
	for key, value := range info {
		switch key {
		case "nickname":
			user.Nickname = value.(string)
		case "introduction":
			user.Introduction = value.(string)
		case "telephone":
			user.Telephone = value.(string)
		case "qq":
			user.QQ = value.(string)
		case "gender":
			user.Gender = value.(string)
		case "email":
			user.Email = value.(string)
		case "birthday":
			user.Birthday = value.(time.Time)
		}
	}

	err = db.Save(&user)
	if err.Error != nil {
		return err.Error
	}

	return nil
}
