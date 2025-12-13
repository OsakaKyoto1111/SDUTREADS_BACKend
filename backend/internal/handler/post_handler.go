package handler

import (
	"net/http"

	"backend/internal/dto"
	"backend/internal/service"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postSvc service.PostService
	fileSvc *service.FileService
}

func NewPostHandler(p service.PostService, f *service.FileService) *PostHandler {
	return &PostHandler{postSvc: p, fileSvc: f}
}

func (h *PostHandler) Create(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	var req dto.CreatePostRequestMultipart
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid form data")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return respondError(c, http.StatusBadRequest, "failed to parse multipart form")
	}

	files := form.File["files"]

	postID, err := h.postSvc.CreatePostWithFiles(userID, req, files)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusCreated, echo.Map{"id": postID})
}

func (h *PostHandler) Get(c echo.Context) error {
	postID, err := parseIDParam(c, "id")
	if err != nil {
		httpErr := err.(*echo.HTTPError)
		return respondError(c, httpErr.Code, httpErr.Message.(string))
	}

	userID, _ := GetUserIDFromContext(c)

	post, err := h.postSvc.GetPost(postID, userID)
	if err != nil {
		return respondError(c, http.StatusNotFound, err.Error())
	}

	return respondJSON(c, http.StatusOK, post)
}

func (h *PostHandler) Update(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	postID, err := parseIDParam(c, "id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid post id")
	}

	var req dto.UpdatePostRequest
	if err := bindJSON(c, &req); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	if err := h.postSvc.UpdatePost(postID, userID, req); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, echo.Map{"message": "updated"})
}

func (h *PostHandler) Delete(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	postID, err := parseIDParam(c, "id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid post id")
	}

	if err := h.postSvc.DeletePost(postID, userID); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, echo.Map{"message": "deleted"})
}

func (h *PostHandler) Like(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	postID, err := parseIDParam(c, "id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid post id")
	}

	if err := h.postSvc.LikePost(postID, userID); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, echo.Map{"message": "liked"})
}

func (h *PostHandler) Unlike(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	postID, err := parseIDParam(c, "id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid post id")
	}

	if err := h.postSvc.UnlikePost(postID, userID); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return respondJSON(c, http.StatusOK, echo.Map{"message": "unliked"})
}

func (h *PostHandler) AddFiles(c echo.Context) error {
	postID, err := parseIDParam(c, "id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid post id")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid form data")
	}

	files := form.File["files"]
	urls, err := h.fileSvc.SaveFiles(postID, files)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	if err := h.postSvc.AddFiles(postID, urls); err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	return respondJSON(c, http.StatusOK, urls)
}

func (h *PostHandler) MyPosts(c echo.Context) error {
	userID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	posts, err := h.postSvc.GetUserPosts(userID, userID)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	return respondJSON(c, http.StatusOK, posts)
}

func (h *PostHandler) UserPosts(c echo.Context) error {
	viewerID, ok := requireAuth(c)
	if !ok {
		return nil
	}

	targetID, err := parseIDParam(c, "id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid user id")
	}

	posts, err := h.postSvc.GetUserPosts(targetID, viewerID)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	return respondJSON(c, http.StatusOK, posts)
}
