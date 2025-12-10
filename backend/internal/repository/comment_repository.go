package repository

import (
	"fmt"

	"backend/internal/model"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Add(comment *model.Comment) error
	GetRootComments(postID uint, limit, offset int) ([]model.Comment, error)
	GetCommentsByPostID(postID uint) ([]model.Comment, error)
	GetReplies(parentID uint) ([]model.Comment, error)
	GetCommentsFlat(postID uint) ([]model.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Add(comment *model.Comment) error {
	if comment == nil {
		return fmt.Errorf("comment is nil")
	}
	if err := r.db.Create(comment).Error; err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

func (r *commentRepository) GetRootComments(postID uint, limit, offset int) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.
		Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("User").
		Preload("Likes").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("get root comments: %w", err)
	}
	return comments, nil
}

func (r *commentRepository) GetCommentsByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.
		Where("post_id = ?", postID).
		Preload("User").
		Preload("Likes").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("get comments by post id: %w", err)
	}
	return comments, nil
}

func (r *commentRepository) GetReplies(parentID uint) ([]model.Comment, error) {
	var replies []model.Comment
	if err := r.db.
		Where("parent_id = ?", parentID).
		Preload("User").
		Preload("Likes").
		Order("created_at ASC").
		Find(&replies).Error; err != nil {
		return nil, fmt.Errorf("get replies: %w", err)
	}
	return replies, nil
}

func (r *commentRepository) GetCommentsFlat(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.
		Where("post_id = ?", postID).
		Preload("User").
		Preload("Likes").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("get comments flat: %w", err)
	}
	return comments, nil
}
