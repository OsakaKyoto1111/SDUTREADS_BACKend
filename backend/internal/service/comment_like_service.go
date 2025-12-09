package service

import (
	"backend/internal/dto"
	"backend/internal/repository"
	"errors"
)

type CommentLikeService interface {
	Like(commentID uint, userID uint) (*dto.CommentLikeResponse, error)
	Unlike(commentID uint, userID uint) (*dto.CommentLikeResponse, error)
}

type commentLikeService struct {
	repo repository.CommentLikeRepository
}

func NewCommentLikeService(repo repository.CommentLikeRepository) CommentLikeService {
	return &commentLikeService{repo: repo}
}

func (s *commentLikeService) Like(commentID uint, userID uint) (*dto.CommentLikeResponse, error) {
	err := s.repo.LikeComment(commentID, userID)
	if err != nil {
		return nil, errors.New("failed to like comment")
	}

	return &dto.CommentLikeResponse{
		CommentID: commentID,
		Liked:     true,
	}, nil
}

func (s *commentLikeService) Unlike(commentID uint, userID uint) (*dto.CommentLikeResponse, error) {
	err := s.repo.UnlikeComment(commentID, userID)
	if err != nil {
		return nil, errors.New("failed to unlike comment")
	}

	return &dto.CommentLikeResponse{
		CommentID: commentID,
		Liked:     false,
	}, nil
}
