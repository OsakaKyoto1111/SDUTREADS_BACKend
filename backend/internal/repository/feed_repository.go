package repository

import (
	"backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type FeedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) *FeedRepository {
	return &FeedRepository{db: db}
}

func (r *FeedRepository) GetFollowingPosts(userID uint, limit int, cursor *time.Time) ([]model.Post, error) {
	var posts []model.Post

	q := r.db.
		Joins("JOIN follows ON follows.following_id = posts.user_id").
		Where("follows.follower_id = ?", userID).
		Preload("User").
		Preload("Files").
		Preload("Likes").
		Preload("Comments").
		Order("posts.created_at DESC").
		Limit(limit)

	if cursor != nil {
		q = q.Where("posts.created_at < ?", *cursor)
	}

	err := q.Find(&posts).Error
	return posts, err
}

func (r *FeedRepository) GetRecommendedPosts(userID uint, limit int) ([]model.Post, error) {
	var posts []model.Post

	err := r.db.
		Raw(`
            SELECT * FROM posts
            WHERE user_id != ?
            ORDER BY RANDOM()
            LIMIT ?
        `, userID, limit).
		Scan(&posts).Error

	return posts, err
}
