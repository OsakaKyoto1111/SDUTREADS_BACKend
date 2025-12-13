package repository

import (
	"fmt"

	"backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostRepository interface {
	CreatePost(post *model.Post) error
	UpdateFields(id uint, fields map[string]interface{}) error
	DeletePost(id uint, userID uint) error
	AddFiles(files []model.File) error
	LikePost(postID, userID uint) error
	UnlikePost(postID, userID uint) error
	FindByID(id uint) (*model.Post, error)
	GetByUser(userID uint) ([]model.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(post *model.Post) error {
	if post == nil {
		return fmt.Errorf("post is nil")
	}
	if err := r.db.Create(post).Error; err != nil {
		return fmt.Errorf("create post: %w", err)
	}
	return nil
}

func (r *postRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	if err := r.db.Model(&model.Post{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return fmt.Errorf("update post fields: %w", err)
	}
	return nil
}

func (r *postRepository) DeletePost(id uint, userID uint) error {
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Post{}).Error; err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}

func (r *postRepository) AddFiles(files []model.File) error {
	if len(files) == 0 {
		return nil
	}
	if err := r.db.Create(&files).Error; err != nil {
		return fmt.Errorf("add files: %w", err)
	}
	return nil
}

func (r *postRepository) LikePost(postID, userID uint) error {
	if postID == 0 || userID == 0 {
		return fmt.Errorf("invalid ids")
	}
	like := model.PostLike{
		PostID: postID,
		UserID: userID,
	}
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&like).Error; err != nil {
		return fmt.Errorf("like post: %w", err)
	}
	return nil
}

func (r *postRepository) UnlikePost(postID, userID uint) error {
	if err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&model.PostLike{}).Error; err != nil {
		return fmt.Errorf("unlike post: %w", err)
	}
	return nil
}

func (r *postRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := r.db.
		Where("id = ?", id).
		Preload("Files").
		Preload("Likes").
		First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find post: %w", err)
	}
	return &post, nil
}

func (r *postRepository) GetByUser(userID uint) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.
		Where("user_id = ?", userID).
		Preload("Files").
		Preload("Likes").
		Preload("Comments").
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("get posts by user: %w", err)
	}
	return posts, nil
}
