package handler

import (
	"net/http"

	"backend/internal/dto"
	"backend/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := bindJSON(c, &req); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}
	resp, err := h.svc.Register(req)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}
	return respondJSON(c, http.StatusOK, resp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := bindJSON(c, &req); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}
	resp, err := h.svc.Login(req)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err.Error())
	}
	return respondJSON(c, http.StatusOK, resp)
}
