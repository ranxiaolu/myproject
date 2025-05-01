package model

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

const JwtSecretKey string = "secret"

type Token struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
