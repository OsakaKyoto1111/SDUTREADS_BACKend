package service

import (
	"fmt"
	"mime/multipart"

	"backend/internal/dto"
	"backend/internal/mapper"
	"backend/internal/model"
	"backend/internal/repository"
)

type PostService interface {
	CreatePost(userID uint, req dto.CreatePostRequest) error
	UpdatePost(postID uint, userID uint, req dto.UpdatePostRequest) error
	DeletePost(postID, userID uint) error
	AddFiles(postID uint, urls []string) error
	LikePost(postID, userID uint) error
	UnlikePost(postID, userID uint) error
	GetPost(postID, userID uint) (*dto.PostWithCommentsResponse, error)
	CreatePostWithFiles(userID uint, req dto.CreatePostRequestMultipart, files []*multipart.FileHeader) (uint, error)
	GetUserPosts(targetUserID, viewerID uint) ([]dto.PostResponse, error)
}

type postService struct {
	repo        repository.PostRepository
	commentSvc  CommentService
	commentTree CommentTreeService
	fileSvc     *FileService
}

func (s *postService) CreatePostWithFiles(userID uint, req dto.CreatePostRequestMultipart, files []*multipart.FileHeader) (uint, error) {
	if userID == 0 {
		return 0, fmt.Errorf("unauthorized")
	}

	post := model.Post{
		UserID:      userID,
		Description: req.Description,
	}

	// создаём сам пост
	if err := s.repo.CreatePost(&post); err != nil {
		return 0, err
	}

	// если файлов нет — return
	if len(files) == 0 {
		return post.ID, nil
	}

	urls, err := s.fileSvc.SaveFiles(post.ID, files)
	if err != nil {
		return 0, err
	}

	// сохраняем в БД
	var dbFiles []model.File
	for _, u := range urls {
		dbFiles = append(dbFiles, model.File{
			PostID: post.ID,
			URL:    u,
		})
	}

	if err := s.repo.AddFiles(dbFiles); err != nil {
		return 0, err
	}

	return post.ID, nil
}

func NewPostService(
	postRepo repository.PostRepository,
	commentSvc CommentService,
	commentTree CommentTreeService,
	fileSvc *FileService,
) PostService {
	return &postService{
		repo:        postRepo,
		commentSvc:  commentSvc,
		commentTree: commentTree,
		fileSvc:     fileSvc,
	}
}

func (s *postService) CreatePost(userID uint, req dto.CreatePostRequest) error {
	if userID == 0 {
		return fmt.Errorf("unauthorized")
	}
	post := model.Post{
		UserID:      userID,
		Description: req.Description,
	}
	return s.repo.CreatePost(&post)
}

func (s *postService) UpdatePost(postID uint, userID uint, req dto.UpdatePostRequest) error {
	if postID == 0 || userID == 0 {
		return fmt.Errorf("invalid ids")
	}
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return fmt.Errorf("find post: %w", err)
	}
	if post.UserID != userID {
		return fmt.Errorf("forbidden")
	}
	updates := map[string]interface{}{}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if len(updates) == 0 {
		return nil
	}
	return s.repo.UpdateFields(postID, updates)
}

func (s *postService) DeletePost(postID uint, userID uint) error {
	if postID == 0 || userID == 0 {
		return fmt.Errorf("invalid ids")
	}
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return fmt.Errorf("find post: %w", err)
	}
	if post.UserID != userID {
		return fmt.Errorf("forbidden")
	}
	return s.repo.DeletePost(postID, userID)
}

func (s *postService) AddFiles(postID uint, urls []string) error {
	if postID == 0 {
		return fmt.Errorf("invalid post id")
	}
	var f []model.File
	for _, url := range urls {
		f = append(f, model.File{
			PostID: postID,
			URL:    url,
		})
	}
	return s.repo.AddFiles(f)
}

func (s *postService) LikePost(postID, userID uint) error {
	if postID == 0 || userID == 0 {
		return fmt.Errorf("invalid ids")
	}
	return s.repo.LikePost(postID, userID)
}

func (s *postService) UnlikePost(postID, userID uint) error {
	if postID == 0 || userID == 0 {
		return fmt.Errorf("invalid ids")
	}
	return s.repo.UnlikePost(postID, userID)
}

func (s *postService) GetPost(postID, userID uint) (*dto.PostWithCommentsResponse, error) {
	if postID == 0 {
		return nil, fmt.Errorf("invalid id")
	}
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return nil, fmt.Errorf("find post: %w", err)
	}

	isLiked := false
	for _, l := range post.Likes {
		if l.UserID == userID {
			isLiked = true
			break
		}
	}

	var files []dto.FileResponse
	for _, f := range post.Files {
		files = append(files, dto.FileResponse{ID: f.ID, URL: f.URL})
	}

	tree, err := s.commentTree.GetCommentTree(postID, userID)
	if err != nil {
		return nil, fmt.Errorf("comment tree: %w", err)
	}

	return &dto.PostWithCommentsResponse{
		ID:          post.ID,
		Description: post.Description,
		Files:       files,
		LikesCount:  len(post.Likes),
		IsLiked:     isLiked,
		Comments:    tree,
	}, nil
}

func (s *postService) GetUserPosts(targetUserID, viewerID uint) ([]dto.PostResponse, error) {
	if targetUserID == 0 {
		return nil, fmt.Errorf("invalid id")
	}
	if viewerID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	posts, err := s.repo.GetByUser(targetUserID)
	if err != nil {
		return nil, err
	}

	return mapper.MapPostsToDTO(posts, viewerID), nil
}
