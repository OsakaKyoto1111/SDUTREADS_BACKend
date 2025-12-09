package handler

import (
	"backend/internal/dto"
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService *service.PostService
	fileService *service.FileService
}

func NewPostHandler(p *service.PostService, f *service.FileService) *PostHandler {
	return &PostHandler{postService: p, fileService: f}
}

func getUserIDFromContext(c echo.Context) (uint, bool) {
	val := c.Get("user_id")
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case int64:
		return uint(v), true
	case float64:
		return uint(v), true
	default:
		return 0, false
	}
}

func parseIDParam(p string) (uint, error) {
	id, err := strconv.ParseUint(p, 10, 64)
	return uint(id), err
}

func (h *PostHandler) Create(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	req := new(dto.CreatePostRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := h.postService.CreatePost(userID, *req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "post created"})
}

func (h *PostHandler) Update(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	postID, err := parseIDParam(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid post id"})
	}

	req := new(dto.UpdatePostRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := h.postService.UpdatePost(postID, userID, *req); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "forbidden"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "post updated"})
}

func (h *PostHandler) Delete(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	postID, err := parseIDParam(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid post id"})
	}

	if err := h.postService.DeletePost(postID, userID); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "forbidden"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "post deleted"})
}

func (h *PostHandler) AddFiles(c echo.Context) error {
	postID, err := parseIDParam(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid post id"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid form"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no files provided"})
	}

	urls, err := h.fileService.SaveFiles(postID, files)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.postService.AddFiles(postID, urls); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"urls": urls})
}

func (h *PostHandler) LikePost(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	postID, err := parseIDParam(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid post id"})
	}

	if err := h.postService.LikePost(postID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "liked"})
}

func (h *PostHandler) UnlikePost(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	postID, err := parseIDParam(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid post id"})
	}

	if err := h.postService.UnlikePost(postID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "unliked"})
}
func (h *PostHandler) GetPost(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	postID, err := parseIDParam(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	resp, err := h.postService.GetPost(postID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}
