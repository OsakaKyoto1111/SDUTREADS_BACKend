package handler

import (
	"backend/internal/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthCheckHandler struct {
	userRepo repository.UserRepository
}

func NewAuthCheckHandler(repo repository.UserRepository) *AuthCheckHandler {
	return &AuthCheckHandler{
		userRepo: repo,
	}
}

func (h *AuthCheckHandler) CheckToken(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token: no user id",
		})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "user does not exist",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "valid",
	})
}
