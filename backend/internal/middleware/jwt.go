package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

func JWT(secret string) echo.MiddlewareFunc {
	if secret == "" {
		log.Fatal("JWT secret is empty!")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("JWT: Missing Authorization header, Path: %s", c.Path())
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Printf("JWT: Invalid Authorization format, Path: %s, Header: %s", c.Path(), authHeader)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			tokenString := strings.TrimSpace(parts[1])
			if tokenString == "" {
				log.Printf("JWT: Empty token, Path: %s", c.Path())
				return echo.NewHTTPError(http.StatusUnauthorized, "empty token")
			}

			tokenParts := strings.Split(tokenString, ".")
			if len(tokenParts) != 3 {
				log.Printf("JWT: Invalid token format (expected 3 parts, got %d), Path: %s, TokenLength: %d",
					len(tokenParts), c.Path(), len(tokenString))
				return echo.NewHTTPError(http.StatusUnauthorized, "malformed token")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid signing method")
				}
				return []byte(secret), nil
			})

			if err != nil {
				if jwtErr, ok := err.(*jwt.ValidationError); ok {
					if jwtErr.Errors&jwt.ValidationErrorExpired != 0 {
						return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
					}
					if jwtErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
						return echo.NewHTTPError(http.StatusUnauthorized, "invalid token signature")
					}
					if jwtErr.Errors&jwt.ValidationErrorMalformed != 0 {
						return echo.NewHTTPError(http.StatusUnauthorized, "malformed token")
					}
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userID, _ := claims["user_id"].(float64)
				nickname, _ := claims["nickname"].(string)
				c.Set("user_id", uint(userID))
				c.Set("user_nickname", nickname)
				return next(c)
			}

			if claims, ok := token.Claims.(*JWTClaims); ok {
				c.Set("user", claims)
				c.Set("user_id", claims.UserID)
				c.Set("user_nickname", claims.Nickname)
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}
	}
}
