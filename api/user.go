package api

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	//"github.com/gin-gonic/gin"
	//"time"
	//"github.com/gin-gonic/gin"
	"myproject/dao"
	"myproject/model"
	"myproject/myutils"
	"myproject/service"
	"net/http"
)

func RegisterUser(ctx context.Context, c *app.RequestContext) {
	// 参数校验和请求转发
	var user model.User
	//收到一个包含 JSON 数据的 POST 请求时，使用 c.BindJSON 方法将请求体中的 JSON 数据解析到 User 结构体中
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, map[string]string{"error": "参数错误" + err.Error()})
		return
	}
	//查询用户是否存在
	exists, err := dao.FindUser(&user)
	if !err {
		c.JSON(500, map[string]string{"error": "查询用户发生错误"}) //400 表示服务器端错误的响应状态码
		return
	}
	if exists {
		c.JSON(409, map[string]string{"error": "用户已存在"})
		return
	}
	//传入数据库
	if err := service.RegisterUser(user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "10000", "info": "success"})
	//return nil
}

func LoginUser(ctx context.Context, c *app.RequestContext) {

	var user model.User
	//收到一个包含 JSON 数据的 POST 请求时，使用 c.BindJSON 方法将请求体中的 JSON 数据解析到 User 结构体中
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "解析失败" + err.Error()})
		return
	}

	if exist, _ := dao.FindUser(&user); !exist {
		c.JSON(http.StatusNotFound, utils.H{"info": "用户不存在或查找错误"})
		return
	}

	refreshToken, token, err := service.CreateToken(user.Username, user.Password)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, map[string]string{"info": "生成JWT令牌失败" + err.Error()})
	//	return
	//}
	if err != nil {
		fmt.Printf("!!! JWT生成错误详情: %+v\n", err) // 添加此行
		c.JSON(http.StatusUnauthorized, map[string]string{"info": "生成JWT令牌失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": 10000, "info": "success", "data": utils.H{"refresh_token": refreshToken, "token": token}})

}

func RefreshUserToken(ctx context.Context, c *app.RequestContext) {

	var myToken model.Token
	if err := c.BindJSON(&myToken); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()}) //400
		return
	}
	//解析token
	claims, err := myutils.ParseToken(myToken.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.H{
			"info": "token无效",
		})
		return
	}

	rep := model.CustomClaims{
		Username: claims.Username,
		Password: claims.Password,
	}
	//生成新token并刷新
	newRefreshToken, newToken, err := service.CreateToken(rep.Username, rep.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()}) //500
		return
	}

	// 返回token和refresh token
	c.JSON(http.StatusOK, utils.H{"status": 10000, "info": "success", "data": utils.H{"refresh_token": newRefreshToken, "token": newToken}})
}

func ChangePassword(ctx context.Context, c *app.RequestContext) {

	//	c.JSON(400, map[string]string{"error": err.Error()})
	//}
	var rep struct {
		Username    string `json:"username" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.BindJSON(&rep); err != nil {
		c.JSON(400, map[string]string{"error": "参数错误" + err.Error()})
		return
	}
	err := service.ChangePassword(rep.Username, rep.OldPassword, rep.NewPassword)
	if err != nil {
		c.JSON(500, map[string]string{"error": "密码更新失败" + err.Error()})
		return
	}
	c.JSON(200, map[string]string{"status": "10000", "info": "success"})
}

func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id") //
	user, err := service.GetUserInfo(userID)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"status": "10000", "info": "success", "data": user})
}

func UpdateUserInfo(ctx context.Context, c *app.RequestContext) {

	var user model.User
	//var info map[string]interface{}
	//token,err := service.LoginUser(user.Username,user.Password)

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	err = service.UpdateUserInfo(&user)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "10000", "info": "success"})
}
