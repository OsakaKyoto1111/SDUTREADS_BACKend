package handler

import (
	"net/http"
	"strconv"

	"backend/internal/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func resolveCurrentUserID(c echo.Context) (uint, error) {
	user := c.Get("user")
	if user == nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "missing jwt token")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "invalid token context")
	}

	claims, ok := token.Claims.(*dto.JwtCustomClaims)
	if !ok || !token.Valid {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}

	if claims.UserID == 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "user id missing in claims")
	}

	return claims.UserID, nil
}

func parseUintParam(raw string) (uint, error) {
	if raw == "" {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "id must be a positive integer")
	}
	return uint(value), nil
}
