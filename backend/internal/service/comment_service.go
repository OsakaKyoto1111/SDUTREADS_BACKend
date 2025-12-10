package service

import (
	"fmt"
	"sort"

	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
)

type CommentService interface {
	AddComment(postID, userID uint, req dto.AddCommentRequest) error
	GetCommentsTree(postID uint, userID uint, page, limit int) (*dto.CommentListResponse, error)
}

type commentService struct {
	repo repository.CommentRepository
	like repository.CommentLikeRepository
}

func NewCommentService(r repository.CommentRepository, l repository.CommentLikeRepository) CommentService {
	return &commentService{repo: r, like: l}
}

func (s *commentService) AddComment(postID, userID uint, req dto.AddCommentRequest) error {
	if postID == 0 || userID == 0 {
		return fmt.Errorf("invalid ids")
	}
	if req.Text == "" {
		return fmt.Errorf("text is required")
	}

	comment := model.Comment{
		PostID:   postID,
		UserID:   userID,
		ParentID: req.ParentID,
		Text:     req.Text,
	}
	return s.repo.Add(&comment)
}

func (s *commentService) GetCommentsTree(postID uint, userID uint, page, limit int) (*dto.CommentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if postID == 0 {
		return nil, fmt.Errorf("invalid post id")
	}

	allComments, err := s.repo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("load comments: %w", err)
	}

	if len(allComments) == 0 {
		return &dto.CommentListResponse{
			Comments: []dto.CommentDTO{},
			Page:     page,
			PageSize: limit,
			Total:    0,
		}, nil
	}

	allIDs := make([]uint, 0, len(allComments))
	for _, c := range allComments {
		allIDs = append(allIDs, c.ID)
	}
	likesMap, err := s.like.LikesCountsForComments(allIDs)
	if err != nil {
		return nil, fmt.Errorf("likes counts: %w", err)
	}
	likedMap, err := s.like.LikedByUser(allIDs, userID)
	if err != nil {
		return nil, fmt.Errorf("liked by user: %w", err)
	}

	type node struct {
		dto dto.CommentDTO
		raw model.Comment
	}
	nodes := make(map[uint]*node, len(allComments))
	var roots []uint

	for _, c := range allComments {
		lc := 0
		if v, ok := likesMap[c.ID]; ok {
			lc = v
		}
		isLiked := false
		if userID != 0 {
			if v, ok := likedMap[c.ID]; ok {
				isLiked = v
			}
		}

		d := dto.CommentDTO{
			ID:       c.ID,
			PostID:   c.PostID,
			UserID:   c.UserID,
			ParentID: c.ParentID,
			Text:     c.Text,
			Likes:    lc,
			IsLiked:  isLiked,
			User: dto.UserShortDTO{
				ID:       c.User.ID,
				Nickname: c.User.Nickname,
				Avatar:   c.User.AvatarURL,
			},
			Replies:   []dto.CommentDTO{},
			CreatedAt: c.CreatedAt,
		}
		nodes[c.ID] = &node{dto: d, raw: c}
		if c.ParentID == nil {
			roots = append(roots, c.ID)
		}
	}

	sort.Slice(roots, func(i, j int) bool {
		return nodes[roots[i]].raw.CreatedAt.After(nodes[roots[j]].raw.CreatedAt)
	})

	for _, n := range nodes {
		if n.dto.ParentID != nil {
			parent, ok := nodes[*n.dto.ParentID]
			if ok {
				if parent.dto.ParentID != nil && *parent.dto.ParentID == n.dto.ID {
					continue
				}
				parent.dto.Replies = append(parent.dto.Replies, n.dto)
			}
		}
	}

	for _, n := range nodes {
		if len(n.dto.Replies) > 1 {
			sort.Slice(n.dto.Replies, func(i, j int) bool {
				return n.dto.Replies[i].CreatedAt.Before(n.dto.Replies[j].CreatedAt)
			})
		}
	}

	totalRoots := len(roots)
	start := (page - 1) * limit
	if start > totalRoots {
		start = totalRoots
	}
	end := start + limit
	if end > totalRoots {
		end = totalRoots
	}
	pagedRoots := roots[start:end]

	res := make([]dto.CommentDTO, 0, len(pagedRoots))
	for _, id := range pagedRoots {
		res = append(res, nodes[id].dto)
	}

	resp := &dto.CommentListResponse{
		Comments: res,
		Page:     page,
		PageSize: limit,
		Total:    int64(totalRoots),
	}
	return resp, nil
}
