package handler

import (
	"net/http"

	"backend/internal/dto"
	"backend/internal/service"
	"backend/internal/utils"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}

	authResp, err := h.authService.Register(req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}

	return utils.Success(c, http.StatusCreated, authResp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}

	authResp, err := h.authService.Login(req)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err)
	}

	return utils.Success(c, http.StatusOK, authResp)
}
