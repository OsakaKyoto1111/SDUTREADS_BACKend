package service

import (
	"backend/internal/dto"
	"backend/internal/repository"
)

type CommentTreeService struct {
	repo     *repository.CommentRepository
	likeRepo repository.CommentLikeRepository
}

func NewCommentTreeService(repo *repository.CommentRepository, likeRepo repository.CommentLikeRepository) *CommentTreeService {
	return &CommentTreeService{repo: repo, likeRepo: likeRepo}
}

func (s *CommentTreeService) GetCommentTree(postID, userID uint) ([]dto.CommentTree, error) {
	comments, err := s.repo.GetCommentsFlat(postID)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
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

	// Build tree
	for _, node := range dtoMap {
		if node.ParentID == nil {
			roots = append(roots, *node)
		} else {
			parent := dtoMap[*node.ParentID]
			if parent != nil {
				parent.Replies = append(parent.Replies, *node)
			}
		}
	}

	return roots, nil
}
