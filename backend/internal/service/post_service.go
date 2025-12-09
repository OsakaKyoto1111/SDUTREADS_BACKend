package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo}
}

func (s *PostService) CreatePost(userID uint, req dto.CreatePostRequest) error {
	post := model.Post{
		UserID:      userID,
		Description: req.Description,
	}
	return s.postRepo.CreatePost(&post)
}

func (s *PostService) UpdatePost(postID uint, userID uint, req dto.UpdatePostRequest) error {
	post := model.Post{
		ID:          postID,
		UserID:      userID,
		Description: req.Description,
	}
	return s.postRepo.UpdatePost(&post)
}

func (s *PostService) DeletePost(postID uint, userID uint) error {
	return s.postRepo.DeletePost(postID, userID)
}

func (s *PostService) AddFiles(postID uint, files []string) error {
	var f []model.File
	for _, url := range files {
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

func (s *PostService) AddComment(postID, userID uint, req dto.AddCommentRequest) error {
	comment := model.Comment{
		PostID:   postID,
		UserID:   userID,
		ParentID: req.ParentID,
		Text:     req.Text,
	}
	return s.postRepo.AddComment(&comment)
}

func (s *PostService) LikeComment(commentID, userID uint) error {
	return s.postRepo.LikeComment(commentID, userID)
}

func (s *PostService) UnlikeComment(commentID, userID uint) error {
	return s.postRepo.UnlikeComment(commentID, userID)
}
