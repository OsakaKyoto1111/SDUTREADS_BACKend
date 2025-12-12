package repository

import (
	"backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type FeedRepository interface {
	GetFollowingPosts(userID uint, limit int, cursor *time.Time) ([]model.Post, error)
	GetRecommendedPosts(userID uint, limit int) ([]model.Post, error)
}

type feedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) FeedRepository {
	return &feedRepository{db: db}
}

func (r *feedRepository) GetFollowingPosts(userID uint, limit int, cursor *time.Time) ([]model.Post, error) {
	var posts []model.Post

	q := r.db.
		Joins("JOIN followers ON followers.user_id = posts.user_id").
		Where("followers.follower_id = ?", userID).
		Preload("User").
		Preload("Files").
		Preload("Likes").
		Preload("Comments").
		Order("posts.created_at DESC").
		Limit(limit)

	if cursor != nil {
		q = q.Where("posts.created_at < ?", *cursor)
	}

	if err := q.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *feedRepository) GetRecommendedPosts(userID uint, limit int) ([]model.Post, error) {
	var posts []model.Post

	// Сначала выберем ID постов raw-запросом
	var ids []uint
	if err := r.db.
		Raw(`
            SELECT id FROM posts
            WHERE user_id != ?
            ORDER BY RANDOM()
            LIMIT ?
        `, userID, limit).
		Scan(&ids).Error; err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []model.Post{}, nil
	}

	// Теперь загрузим сами посты с Preload
	if err := r.db.
		Where("id IN ?", ids).
		Preload("User").
		Preload("Files").
		Preload("Likes").
		Preload("Comments").
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
