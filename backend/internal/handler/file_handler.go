package handler

import (
	"net/http"

	"backend/internal/service"

	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	fileSvc *service.FileService
}

func NewFileHandler(s *service.FileService) *FileHandler {
	return &FileHandler{fileSvc: s}
}

func (h *FileHandler) Upload(c echo.Context) error {
	postID, ok := parseIDParam(c, "id")
	if ok != nil {
		postID = 0 // for avatar upload, no post ID
	}

	form, err := c.MultipartForm()
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	files := form.File["files"]

	urls, err := h.fileSvc.SaveFiles(postID, files)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, urls)
}
