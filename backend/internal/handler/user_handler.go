package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"backend/internal/dto"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/internal/utils"

	"github.com/labstack/echo/v4"
)

// UserHandler exposes user profile endpoints.
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler builds a UserHandler.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe returns the currently authenticated user.
func (h *UserHandler) GetMe(c echo.Context) error {
	userID, err := resolveCurrentUserID(c)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err)
	}
	user, err := h.userService.GetMe(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err)
	}
	return utils.Success(c, http.StatusOK, user)
}

// UpdateProfile patches the authenticated user's profile.
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID, err := resolveCurrentUserID(c)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err)
	}
	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}
	user, err := h.userService.UpdateProfile(c.Request().Context(), userID, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}
	return utils.Success(c, http.StatusOK, user)
}

// SearchUsers returns profiles matching the query.
func (h *UserHandler) SearchUsers(c echo.Context) error {
	query := c.QueryParam("query")
	if strings.TrimSpace(query) == "" {
		return utils.Error(c, http.StatusBadRequest, errors.New("query is required"))
	}

	limit := 25
	if raw := c.QueryParam("limit"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			return utils.Error(c, http.StatusBadRequest, errors.New("limit must be a positive integer"))
		}
		limit = value
	}

	users, err := h.userService.SearchUsers(c.Request().Context(), query, limit)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err)
	}
	return utils.Success(c, http.StatusOK, users)
}

// GetByID returns another user's profile.
func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err)
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return utils.Error(c, http.StatusNotFound, err)
		}
		return utils.Error(c, http.StatusInternalServerError, err)
	}
	return utils.Success(c, http.StatusOK, user)
}
