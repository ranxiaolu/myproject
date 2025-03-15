package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dgrijalva/jwt-go"
	"myproject/model"
	"net/http"
	"strings"
)

// AuthMiddleware 是一个验证 Token 的中间件
func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取 Authorization 请求头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(consts.StatusUnauthorized, map[string]string{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查 Token 格式是否为 "Bearer <token>" 两部分
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" { //验证前缀是否为Bearer
			c.JSON(consts.StatusUnauthorized, map[string]string{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// 获取 Token 并进行验证
		tokenString := parts[1]
		//
		claims := jwt.MapClaims{}
		//解析JWT字符串															//提取密匙
		parsedToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(model.JwtSecretKey), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.H{"status": 10001, "info": "invalid token"})
			c.Abort()
			return
		}
		//验证有效性
		if !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, utils.H{"status": 10001, "info": "invalid token"})
			c.Abort()
			return
		}

		// 如果 Token 验证通过，提取用户信息并继续处理请求
		c.Set("UserID", claims["ID"])
		c.Set("Username", claims["Username"])
		c.Next(ctx)
	}
}
