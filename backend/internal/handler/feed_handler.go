package handler

import (
	"backend/internal/service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type FeedHandler struct {
	service *service.FeedService
}

func NewFeedHandler(s *service.FeedService) *FeedHandler {
	return &FeedHandler{s}
}

func (h *FeedHandler) GetFeed(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	limit := 20
	cursorParam := c.QueryParam("cursor")

	var cursor *time.Time = nil
	if cursorParam != "" {
		t, err := time.Parse(time.RFC3339, cursorParam)
		if err == nil {
			cursor = &t
		}
	}

	resp, err := h.service.GetFeed(userID, limit, cursor)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}
