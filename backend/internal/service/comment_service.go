package service

import (
	"backend/internal/dto"
	"backend/internal/mapper"
	"backend/internal/model"
	"backend/internal/repository"
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

// recursively load replies
func (s *CommentService) loadRepliesRecursive(comment model.Comment, userID uint) (dto.CommentDTO, error) {
	isLiked, _ := s.likeRepo.IsLiked(comment.ID, userID)
	likesCount, _ := s.likeRepo.LikesCount(comment.ID)

	dtoComment := mapper.MapCommentToDTO(comment, isLiked, likesCount)

	replies, _ := s.commentRepo.GetReplies(comment.ID)
	dtoComment.Replies = []dto.CommentDTO{}

	for _, r := range replies {
		child, _ := s.loadRepliesRecursive(r, userID)
		dtoComment.Replies = append(dtoComment.Replies, child)
	}

	return dtoComment, nil
}

func (s *CommentService) GetComments(postID uint, userID uint, limit, offset int) ([]dto.CommentDTO, error) {
	roots, err := s.commentRepo.GetRootComments(postID, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []dto.CommentDTO
	for _, c := range roots {
		dtoComment, _ := s.loadRepliesRecursive(c, userID)
		result = append(result, dtoComment)
	}
	return result, nil
}
