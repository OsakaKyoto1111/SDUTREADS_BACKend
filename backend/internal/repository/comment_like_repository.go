package repository

import (
	"backend/internal/model"

	"gorm.io/gorm"
)

type CommentLikeRepository interface {
	LikeComment(commentID uint, userID uint) error
	UnlikeComment(commentID uint, userID uint) error
	IsLiked(commentID uint, userID uint) (bool, error)
	LikesCount(commentID uint) (int, error)
}

type commentLikeRepository struct {
	db *gorm.DB
}

func NewCommentLikeRepository(db *gorm.DB) CommentLikeRepository {
	return &commentLikeRepository{db: db}
}

func (r *commentLikeRepository) LikeComment(commentID uint, userID uint) error {
	like := model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}

	return r.db.
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		FirstOrCreate(&like).Error
}

func (r *commentLikeRepository) UnlikeComment(commentID uint, userID uint) error {
	return r.db.
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Delete(&model.CommentLike{}).Error
}

func (r *commentLikeRepository) IsLiked(commentID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *commentLikeRepository) LikesCount(commentID uint) (int, error) {
	var count int64
	err := r.db.Model(&model.CommentLike{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error
	return int(count), err
}
