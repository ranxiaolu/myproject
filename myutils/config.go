package myutils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"myproject/model"
	"time"
)

func ParseToken(refreshToken string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(model.JwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.CustomClaims)
	if ok && token.Valid {
		return claims, nil
	}
	//验证有效期
	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, fmt.Errorf("token is expired")
	}

	return nil, err
}
