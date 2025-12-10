package service

import (
	"fmt"

	"backend/internal/dto"
	"backend/internal/repository"
)

type CommentTreeService interface {
	GetCommentTree(postID, userID uint) ([]dto.CommentTree, error)
}

type commentTreeService struct {
	repo     repository.CommentRepository
	likeRepo repository.CommentLikeRepository
}

func NewCommentTreeService(repo repository.CommentRepository, likeRepo repository.CommentLikeRepository) CommentTreeService {
	return &commentTreeService{repo: repo, likeRepo: likeRepo}
}

func (s *commentTreeService) GetCommentTree(postID, userID uint) ([]dto.CommentTree, error) {
	if postID == 0 {
		return nil, fmt.Errorf("invalid post id")
	}

	comments, err := s.repo.GetCommentsFlat(postID)
	if err != nil {
		return nil, fmt.Errorf("load comments flat: %w", err)
	}

	dtoMap := map[uint]*dto.CommentTree{}
	var roots []dto.CommentTree

	for _, c := range comments {
		isLiked, _ := s.likeRepo.IsLiked(c.ID, userID)

		node := &dto.CommentTree{
			ID:        c.ID,
			PostID:    c.PostID,
			UserID:    c.UserID,
			ParentID:  c.ParentID,
			Text:      c.Text,
			Likes:     len(c.Likes),
			IsLiked:   isLiked,
			CreatedAt: c.CreatedAt,
			User: dto.UserShortDTO{
				ID:       c.User.ID,
				Nickname: c.User.Nickname,
				Avatar:   c.User.AvatarURL,
			},
			Replies: []dto.CommentTree{},
		}

		dtoMap[c.ID] = node
	}

	for id, node := range dtoMap {
		if node.ParentID == nil {
			roots = append(roots, *node)
			continue
		}
		parent, ok := dtoMap[*node.ParentID]
		if !ok {
			roots = append(roots, *node)
			continue
		}
		if parent.ParentID != nil && *parent.ParentID == id {
			roots = append(roots, *node)
			continue
		}
		parent.Replies = append(parent.Replies, *node)
	}

	return roots, nil
}
