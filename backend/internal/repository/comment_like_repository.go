package repository

import (
	"fmt"

	"backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentLikeRepository interface {
	LikeComment(commentID uint, userID uint) error
	UnlikeComment(commentID uint, userID uint) error
	IsLiked(commentID uint, userID uint) (bool, error)
	LikesCount(commentID uint) (int, error)

	LikesCountsForComments(commentIDs []uint) (map[uint]int, error)
	LikedByUser(commentIDs []uint, userID uint) (map[uint]bool, error)
}

type commentLikeRepository struct {
	db *gorm.DB
}

func NewCommentLikeRepository(db *gorm.DB) CommentLikeRepository {
	return &commentLikeRepository{db: db}
}

func (r *commentLikeRepository) LikeComment(commentID uint, userID uint) error {
	if commentID == 0 || userID == 0 {
		return fmt.Errorf("invalid ids")
	}

	like := model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}

	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&like).Error; err != nil {
		return fmt.Errorf("create comment_like: %w", err)
	}
	return nil
}

func (r *commentLikeRepository) UnlikeComment(commentID uint, userID uint) error {
	if commentID == 0 || userID == 0 {
		return fmt.Errorf("invalid ids")
	}
	if err := r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).
		Delete(&model.CommentLike{}).Error; err != nil {
		return fmt.Errorf("delete comment_like: %w", err)
	}
	return nil
}

func (r *commentLikeRepository) IsLiked(commentID uint, userID uint) (bool, error) {
	if commentID == 0 || userID == 0 {
		return false, nil
	}
	var cnt int64
	if err := r.db.Model(&model.CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&cnt).Error; err != nil {
		return false, fmt.Errorf("count is_liked: %w", err)
	}
	return cnt > 0, nil
}

func (r *commentLikeRepository) LikesCount(commentID uint) (int, error) {
	if commentID == 0 {
		return 0, nil
	}
	var cnt int64
	if err := r.db.Model(&model.CommentLike{}).
		Where("comment_id = ?", commentID).
		Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("count likes: %w", err)
	}
	return int(cnt), nil
}

func (r *commentLikeRepository) LikesCountsForComments(commentIDs []uint) (map[uint]int, error) {
	result := make(map[uint]int)
	if len(commentIDs) == 0 {
		return result, nil
	}

	type row struct {
		CommentID uint
		Count     int64
	}

	var rows []row
	if err := r.db.Model(&model.CommentLike{}).
		Select("comment_id, count(*) as count").
		Where("comment_id IN ?", commentIDs).
		Group("comment_id").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("likes counts for comments: %w", err)
	}

	for _, r := range rows {
		result[r.CommentID] = int(r.Count)
	}
	return result, nil
}

func (r *commentLikeRepository) LikedByUser(commentIDs []uint, userID uint) (map[uint]bool, error) {
	res := make(map[uint]bool)
	if len(commentIDs) == 0 || userID == 0 {
		return res, nil
	}

	var likes []model.CommentLike
	if err := r.db.
		Select("comment_id").
		Where("user_id = ? AND comment_id IN ?", userID, commentIDs).
		Find(&likes).Error; err != nil {
		return nil, fmt.Errorf("liked by user query: %w", err)
	}

	for _, l := range likes {
		res[l.CommentID] = true
	}
	return res, nil
}
