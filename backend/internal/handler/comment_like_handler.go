package handler

import (
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CommentLikeHandler struct {
	service service.CommentLikeService
}

func NewCommentLikeHandler(service service.CommentLikeService) *CommentLikeHandler {
	return &CommentLikeHandler{service: service}
}

func (h *CommentLikeHandler) Like(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid comment id"})
	}

	resp, err := h.service.Like(uint(commentID), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *CommentLikeHandler) Unlike(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid comment id"})
	}

	resp, err := h.service.Unlike(uint(commentID), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}
