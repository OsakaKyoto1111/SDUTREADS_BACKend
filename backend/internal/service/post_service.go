package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
	"errors"
)

type PostService struct {
	postRepo    *repository.PostRepository
	commentSvc  *CommentService
	commentTree *CommentTreeService
}

func NewPostService(postRepo *repository.PostRepository, commentSvc *CommentService, commentTree *CommentTreeService) *PostService {
	return &PostService{postRepo: postRepo, commentSvc: commentSvc, commentTree: commentTree}
}

func (s *PostService) CreatePost(userID uint, req dto.CreatePostRequest) error {
	post := model.Post{
		UserID:      userID,
		Description: req.Description,
	}
	return s.postRepo.CreatePost(&post)
}

func (s *PostService) UpdatePost(postID uint, userID uint, req dto.UpdatePostRequest) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.New("forbidden")
	}
	updates := map[string]interface{}{}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if len(updates) == 0 {
		return nil
	}
	return s.postRepo.UpdateFields(postID, updates)
}

func (s *PostService) DeletePost(postID uint, userID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.New("forbidden")
	}
	return s.postRepo.DeletePost(postID, userID)
}

func (s *PostService) AddFiles(postID uint, urls []string) error {
	var f []model.File
	for _, url := range urls {
		f = append(f, model.File{
			PostID: postID,
			URL:    url,
		})
	}
	return s.postRepo.AddFiles(f)
}

func (s *PostService) LikePost(postID, userID uint) error {
	return s.postRepo.LikePost(postID, userID)
}

func (s *PostService) UnlikePost(postID, userID uint) error {
	return s.postRepo.UnlikePost(postID, userID)
}

func (s *PostService) GetPost(postID, userID uint) (*dto.PostWithCommentsResponse, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
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
		return nil, err
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
