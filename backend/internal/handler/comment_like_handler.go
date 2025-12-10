package handler

import (
	"net/http"

	"backend/internal/service"

	"github.com/labstack/echo/v4"
)

type CommentLikeHandler struct {
	svc service.CommentLikeService
}

func NewCommentLikeHandler(s service.CommentLikeService) *CommentLikeHandler {
	return &CommentLikeHandler{svc: s}
}

func (h *CommentLikeHandler) Like(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	commentID, err := parseIDParam(c, "comment_id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid id")
	}

	resp, err := h.svc.Like(commentID, userID)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, resp)
}

func (h *CommentLikeHandler) Unlike(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	commentID, err := parseIDParam(c, "comment_id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid id")
	}

	resp, err := h.svc.Unlike(commentID, userID)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, resp)
}
