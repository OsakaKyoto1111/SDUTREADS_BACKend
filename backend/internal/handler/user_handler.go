package handler

import (
	"net/http"

	"backend/internal/dto"
	"backend/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	idStr := c.Param("id")

	if idStr == "" {
		userID, ok := GetUserIDFromContext(c)
		if !ok {
			return respondError(c, http.StatusUnauthorized, "unauthorized")
		}
		resp, err := h.svc.GetUser(userID)
		if err != nil {
			return respondError(c, http.StatusInternalServerError, err.Error())
		}
		return respondJSON(c, http.StatusOK, resp)
	}

	id, err := parseIDParam(c, "id")
	if err != nil {
		httpErr := err.(*echo.HTTPError)
		return respondError(c, httpErr.Code, httpErr.Message.(string))
	}

	resp, err := h.svc.GetUser(id)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return respondJSON(c, http.StatusOK, resp)
}

func (h *UserHandler) Update(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	var body dto.UpdateUserDTO
	if err := bindJSON(c, &body); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	updated, err := h.svc.UpdateUser(userID, body)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return respondJSON(c, http.StatusOK, updated)
}

func (h *UserHandler) Delete(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	if err := h.svc.DeleteUser(userID); err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	return respondJSON(c, http.StatusOK, echo.Map{"message": "user deleted"})
}

func (h *UserHandler) Search(c echo.Context) error {
	q := c.QueryParam("q")
	users, err := h.svc.SearchUsersWithCounts(q)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return respondJSON(c, http.StatusOK, users)
}

func (h *UserHandler) Follow(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	targetID, err := parseIDParam(c, "id")
	if err != nil {
		httpErr := err.(*echo.HTTPError)
		return respondError(c, httpErr.Code, httpErr.Message.(string))
	}

	if err := h.svc.Follow(userID, targetID); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, echo.Map{"message": "followed"})
}

func (h *UserHandler) Unfollow(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	targetID, err := parseIDParam(c, "id")
	if err != nil {
		httpErr := err.(*echo.HTTPError)
		return respondError(c, httpErr.Code, httpErr.Message.(string))
	}

	if err := h.svc.Unfollow(userID, targetID); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, echo.Map{"message": "unfollowed"})
}
