package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
	"sort"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	likeRepo    repository.CommentLikeRepository
}

func NewCommentService(repo *repository.CommentRepository, likeRepo repository.CommentLikeRepository) *CommentService {
	return &CommentService{commentRepo: repo, likeRepo: likeRepo}
}

func (s *CommentService) AddComment(postID, userID uint, req dto.AddCommentRequest) error {
	comment := model.Comment{
		PostID:   postID,
		UserID:   userID,
		ParentID: req.ParentID,
		Text:     req.Text,
	}
	return s.commentRepo.Add(&comment)
}

// GetCommentsTree - returns paginated root comments and their entire reply trees.
// page is 1-based. limit applies to root comments.
func (s *CommentService) GetCommentsTree(postID uint, userID uint, page, limit int) (*dto.CommentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	// 1) load all comments for the post
	allComments, err := s.commentRepo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	// early exit
	if len(allComments) == 0 {
		return &dto.CommentListResponse{
			Comments: []dto.CommentDTO{},
			Page:     page,
			PageSize: limit,
			Total:    0,
		}, nil
	}

	// 2) prepare ids and mapping
	allIDs := make([]uint, 0, len(allComments))
	for _, c := range allComments {
		allIDs = append(allIDs, c.ID)
	}
	likesMap, err := s.likeRepo.LikesCountsForComments(allIDs)
	if err != nil {
		return nil, err
	}
	likedMap, err := s.likeRepo.LikedByUser(allIDs, userID)
	if err != nil {
		return nil, err
	}

	// node storage
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

	// sort roots newest first
	sort.Slice(roots, func(i, j int) bool {
		return nodes[roots[i]].raw.CreatedAt.After(nodes[roots[j]].raw.CreatedAt)
	})

	// build tree
	for _, n := range nodes {
		if n.dto.ParentID != nil {
			parent, ok := nodes[*n.dto.ParentID]
			if ok {
				parent.dto.Replies = append(parent.dto.Replies, n.dto)
			}
		}
	}

	// sort replies by created_at ASC
	for _, n := range nodes {
		if len(n.dto.Replies) > 1 {
			sort.Slice(n.dto.Replies, func(i, j int) bool {
				return n.dto.Replies[i].CreatedAt.Before(n.dto.Replies[j].CreatedAt)
			})
		}
	}

	// paginate roots
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
