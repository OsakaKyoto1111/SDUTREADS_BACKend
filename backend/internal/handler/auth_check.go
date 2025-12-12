package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthCheckHandler struct{}

func NewAuthCheckHandler() *AuthCheckHandler {
	return &AuthCheckHandler{}
}

func (h *AuthCheckHandler) CheckToken(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "valid",
	})
}
