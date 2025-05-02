package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

var jwtKey = []byte("your_secret_key")

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

		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.H{"status": 10001, "info": "invalid token"})
			c.Abort()
			return
		}
		if !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, utils.H{"status": 10001, "info": "invalid token"})
			c.Abort()
			return
		}

		// 如果 Token 验证通过，继续处理请求
		c.Next(ctx)
	}
}
