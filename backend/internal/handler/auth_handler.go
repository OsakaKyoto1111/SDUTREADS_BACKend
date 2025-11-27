package handler

import (
	"net/http"

	"backend/internal/dto"
	"backend/internal/service"
	"backend/internal/utils"

	"github.com/labstack/echo/v4"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler builds an AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register creates a new user account.
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}
	authResp, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}
	return utils.Success(c, http.StatusCreated, authResp)
}

// Login authenticates a user and returns a token.
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}
	authResp, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err)
	}
	return utils.Success(c, http.StatusOK, authResp)
}
