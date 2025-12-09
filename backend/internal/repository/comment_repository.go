package repository

import (
	"backend/internal/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db}
}

func (r *CommentRepository) Add(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// GetRootComments - only root comments (parent_id IS NULL) with User preloaded and optionally Likes preloaded
func (r *CommentRepository) GetRootComments(postID uint, limit, offset int) ([]model.Comment, error) {
	var comments []model.Comment

	err := r.db.
		Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("User").
		Preload("Likes").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error

	return comments, err
}

// GetCommentsByPostID - load all comments for a post with User and Likes preloaded
func (r *CommentRepository) GetCommentsByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment

	err := r.db.
		Where("post_id = ?", postID).
		Preload("User").
		Preload("Likes").
		Order("created_at ASC").
		Find(&comments).Error

	return comments, err
}

func (r *CommentRepository) GetReplies(parentID uint) ([]model.Comment, error) {
	var replies []model.Comment

	err := r.db.
		Where("parent_id = ?", parentID).
		Preload("User").
		Preload("Likes").
		Order("created_at ASC").
		Find(&replies).Error

	return replies, err
}
func (r *CommentRepository) GetCommentsFlat(postID uint) ([]model.Comment, error) {
	var comments []model.Comment

	err := r.db.
		Where("post_id = ?", postID).
		Preload("User").
		Preload("Likes").
		Order("created_at ASC").
		Find(&comments).Error

	return comments, err
}
