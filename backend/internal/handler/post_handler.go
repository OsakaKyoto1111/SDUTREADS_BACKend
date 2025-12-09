package handler

import (
	"backend/internal/dto"
	"backend/internal/service"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(s *service.PostService) *PostHandler {
	return &PostHandler{s}
}

func (h *PostHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	req := new(dto.CreatePostRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	err := h.postService.CreatePost(userID, *req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "post created"})
}

func (h *PostHandler) Update(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	postID := parseID(c.Param("id"))

	req := new(dto.UpdatePostRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	err := h.postService.UpdatePost(postID, userID, *req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "post updated"})
}

func (h *PostHandler) Delete(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	postID := parseID(c.Param("id"))

	err := h.postService.DeletePost(postID, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "post deleted"})
}

func (h *PostHandler) AddFiles(c echo.Context) error {
	postID := parseID(c.Param("id"))

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid form"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no files provided"})
	}

	uploadDir := filepath.Join("uploads", "posts", fmt.Sprint(postID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	var urls []string

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		defer src.Close()

		ext := strings.ToLower(filepath.Ext(file.Filename))
		newName := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir, newName)

		dst, err := os.Create(filePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		url := "/uploads/posts/" + fmt.Sprint(postID) + "/" + newName
		urls = append(urls, url)
	}

	if err := h.postService.AddFiles(postID, urls); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"urls": urls})
}

func (h *PostHandler) AddComment(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	postID := parseID(c.Param("id"))

	req := new(dto.AddCommentRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	err := h.postService.AddComment(postID, userID, *req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, "comment added")
}

func (h *PostHandler) LikePost(c echo.Context) error {
	return h.likeOrUnlike(c, true)
}

func (h *PostHandler) UnlikePost(c echo.Context) error {
	return h.likeOrUnlike(c, false)
}

func (h *PostHandler) likeOrUnlike(c echo.Context, like bool) error {
	userID := c.Get("user_id").(uint)
	postID := parseID(c.Param("id"))

	if like {
		return h.postService.LikePost(postID, userID)
	} else {
		return h.postService.UnlikePost(postID, userID)
	}
}
