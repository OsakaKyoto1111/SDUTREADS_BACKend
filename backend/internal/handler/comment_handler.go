package handler

import (
	"backend/internal/dto"
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{s}
}

func (h *CommentHandler) Add(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	postID, _ := strconv.Atoi(c.Param("id"))

	req := new(dto.AddCommentRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := h.service.AddComment(uint(postID), userID, *req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "comment added"})
}

func (h *CommentHandler) Get(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	postID, _ := strconv.Atoi(c.Param("id"))

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if limit == 0 {
		limit = 20
	}

	comments, err := h.service.GetComments(uint(postID), userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, comments)
}
