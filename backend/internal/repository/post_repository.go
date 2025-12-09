package repository

import (
	"backend/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) UpdatePost(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) DeletePost(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Post{}).Error
}

func (r *PostRepository) AddFiles(files []model.File) error {
	return r.db.Create(&files).Error
}

func (r *PostRepository) LikePost(postID, userID uint) error {
	like := model.PostLike{
		PostID: postID,
		UserID: userID,
	}
	return r.db.FirstOrCreate(&like, like).Error
}

func (r *PostRepository) UnlikePost(postID, userID uint) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&model.PostLike{}).Error
}

func (r *PostRepository) AddComment(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *PostRepository) LikeComment(commentID, userID uint) error {
	like := model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}
	return r.db.FirstOrCreate(&like, like).Error
}

func (r *PostRepository) UnlikeComment(commentID, userID uint) error {
	return r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).
		Delete(&model.CommentLike{}).Error
}
