package repository

import (
	"backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type FeedRepository interface {
	GetFollowingPosts(userID uint, limit int, cursor *time.Time) ([]model.Post, error)
	GetRecommendedPosts(userID uint, limit int, excludeIDs []uint) ([]model.Post, error)
}

type feedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) FeedRepository {
	return &feedRepository{db: db}
}

func (r *feedRepository) GetFollowingPosts(
	userID uint,
	limit int,
	cursor *time.Time,
) ([]model.Post, error) {

	var posts []model.Post

	q := r.db.
		// ✅ PostgreSQL DISTINCT ON
		Select("DISTINCT ON (posts.id) posts.*").
		Joins("JOIN followers ON followers.user_id = posts.user_id").
		Where("followers.follower_id = ?", userID).
		Preload("User").
		Preload("Files").
		Preload("Likes").
		Preload("Comments").
		Order("posts.id, posts.created_at DESC").
		Limit(limit)

	if cursor != nil {
		q = q.Where("posts.created_at < ?", *cursor)
	}

	if err := q.Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *feedRepository) GetRecommendedPosts(userID uint, limit int, excludeIDs []uint) ([]model.Post, error) {
	var posts []model.Post
	var ids []uint

	// ✅ выбираем id рекомендованных постов, исключая уже попавшие в following
	q := r.db.
		Table("posts").
		Select("id").
		Where("user_id != ?", userID)

	if len(excludeIDs) > 0 {
		q = q.Where("id NOT IN ?", excludeIDs)
	}

	if err := q.
		Order("RANDOM()").
		Limit(limit).
		Scan(&ids).Error; err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []model.Post{}, nil
	}

	// ✅ подгружаем сами посты с preload
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
