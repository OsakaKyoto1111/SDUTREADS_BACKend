package middleware

import (
	"backend/internal/dto"

	"github.com/golang-jwt/jwt/v5"
	echjwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWT creates a middleware that validates bearer tokens using the provided secret.
func JWT(secret string) echo.MiddlewareFunc {
	return echjwt.WithConfig(echjwt.Config{
		SigningKey:    []byte(secret),
		SigningMethod: echjwt.AlgorithmHS256,
		TokenLookup:   "header:Authorization:Bearer ",
		ContextKey:    "user",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &dto.JwtCustomClaims{}
		},
	})
}
