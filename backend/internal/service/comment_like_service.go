package service

import (
	"fmt"

	"backend/internal/dto"
	"backend/internal/repository"
)

type CommentLikeService interface {
	Like(commentID, userID uint) (*dto.CommentLikeResponse, error)
	Unlike(commentID, userID uint) (*dto.CommentLikeResponse, error)
}

type commentLikeService struct {
	repo repository.CommentLikeRepository
}

func NewCommentLikeService(r repository.CommentLikeRepository) CommentLikeService {
	return &commentLikeService{repo: r}
}

func (s *commentLikeService) Like(commentID uint, userID uint) (*dto.CommentLikeResponse, error) {
	if commentID == 0 || userID == 0 {
		return nil, fmt.Errorf("invalid ids")
	}
	if err := s.repo.LikeComment(commentID, userID); err != nil {
		return nil, fmt.Errorf("like comment: %w", err)
	}
	return &dto.CommentLikeResponse{CommentID: commentID, Liked: true}, nil
}

func (s *commentLikeService) Unlike(commentID uint, userID uint) (*dto.CommentLikeResponse, error) {
	if commentID == 0 || userID == 0 {
		return nil, fmt.Errorf("invalid ids")
	}
	if err := s.repo.UnlikeComment(commentID, userID); err != nil {
		return nil, fmt.Errorf("unlike comment: %w", err)
	}
	return &dto.CommentLikeResponse{CommentID: commentID, Liked: false}, nil
}
