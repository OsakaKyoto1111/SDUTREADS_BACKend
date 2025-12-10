package dto

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}
