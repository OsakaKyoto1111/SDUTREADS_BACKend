package handler

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/service"

	"github.com/labstack/echo/v4"
)

type FeedHandler struct {
	svc service.FeedService
}

func NewFeedHandler(s service.FeedService) *FeedHandler {
	return &FeedHandler{svc: s}
}

func (h *FeedHandler) Get(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	cursorStr := c.QueryParam("cursor")
	var cursor *time.Time
	if cursorStr != "" {
		t, err := time.Parse(time.RFC3339, cursorStr)
		if err == nil {
			cursor = &t
		}
	}

	resp, err := h.svc.GetFeed(userID, limit, cursor)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	return respondJSON(c, http.StatusOK, resp)
}
