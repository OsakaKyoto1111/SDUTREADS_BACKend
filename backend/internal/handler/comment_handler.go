package handler

import (
	"net/http"

	"backend/internal/dto"
	"backend/internal/service"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	svc service.CommentService
}

func NewCommentHandler(s service.CommentService) *CommentHandler {
	return &CommentHandler{svc: s}
}

func (h *CommentHandler) Add(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	postID, err := parseIDParam(c, "id")
	if err != nil {
		httpErr := err.(*echo.HTTPError)
		return respondError(c, httpErr.Code, httpErr.Message.(string))
	}

	var req dto.AddCommentRequest
	if err := bindJSON(c, &req); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	if err := h.svc.AddComment(postID, userID, req); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, echo.Map{"message": "comment added"})
}

func (h *CommentHandler) GetTree(c echo.Context) error {
	userID, _ := GetUserIDFromContext(c)

	postID, err := parseIDParam(c, "id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid id")
	}

	page := 1
	limit := 10

	tree, err := h.svc.GetCommentsTree(postID, userID, page, limit)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	return respondJSON(c, http.StatusOK, tree)
}
