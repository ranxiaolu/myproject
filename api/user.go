package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"myproject/model"
	"myproject/service"
	"net/http"
)

func RegisterUser(ctx context.Context, c *app.RequestContext) {
	// 参数校验和请求转发
	var password string
	var user model.User
	//收到一个包含 JSON 数据的 POST 请求时，使用 c.BindJSON 方法将请求体中的 JSON 数据解析到 User 结构体中
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	//注册
	if err := service.RegisterUser(user, password); err != nil {
		c.JSON(500, map[string]string{"error": err.Error()}) //表示服务器端错误的响应状态码，
		// 比如服务端在处理请求时遇到了一些错误。
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "10000", "info": "success"})
	//return nil
}

func LoginUser(ctx context.Context, c *app.RequestContext) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//收到一个包含 JSON 数据的 POST 请求时，使用 c.BindJSON 方法将请求体中的 JSON 数据解析到 User 结构体中
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		return
	}
	//获取token
	token, err := service.LoginUser(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	//刷新token获取refresh_token
	refreshToken, err := service.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	c.JSON(http.StatusOK, map[string]interface{}{"status": 10000, "info": "success", "data": utils.H{"refresh_token": refreshToken, "token": token}})

}

func RefreshUserToken(ctx context.Context, c *app.RequestContext) {
	var refreshToken struct {
		Token string `json:"token" binding:"required"`
	}
	//错误饷应
	if err := c.BindJSON(&refreshToken); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()}) //400
		return
	}
	//生成新token
	newToken, err := service.RefreshToken(refreshToken.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()}) //500
		return
	}
	//生成新的refresh_token
	newRefreshToken, err := service.RefreshToken(newToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	// 返回token和refresh token
	c.JSON(http.StatusOK, utils.H{"status": 10000, "info": "success", "data": utils.H{"refresh_token": newRefreshToken, "token": newToken}})
}

func ChangePassword(ctx context.Context, c *app.RequestContext) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
	}
	var newPassword struct {
		NewPassword string `json:"new_password"`
	}
	err := service.ChangePassword(user, user.Password, newPassword.NewPassword)
	if err != nil {

	}
	c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
}

func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	user, err := service.GetUserInfo(userID)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"status": "10000", "info": "success", "data": user})
}

func UpdateUserInfo(ctx context.Context, c *app.RequestContext) {

	var user model.User
	var info map[string]interface{}
	//token,err := service.LoginUser(user.Username,user.Password)

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	err = service.UpdateUserInfo(user, info)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "10000", "info": "success"})
}
