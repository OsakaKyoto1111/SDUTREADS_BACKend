package middleware

import (
	"backend/internal/dto"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echjwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWT(secret string) echo.MiddlewareFunc {
	return echjwt.WithConfig(echjwt.Config{
		SigningKey:    []byte(secret),
		SigningMethod: echjwt.AlgorithmHS256,
		TokenLookup:   "header:Authorization",
		ContextKey:    "user",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &dto.JwtCustomClaims{}
		},
		SuccessHandler: func(c echo.Context) {
			if claims, ok := c.Get("user").(*dto.JwtCustomClaims); ok {
				c.Set("user_id", claims.UserID)
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("invalid or expired token"))
		},
	})
}
