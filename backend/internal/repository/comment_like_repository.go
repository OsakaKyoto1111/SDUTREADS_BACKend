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
	err := r.db.Model(&model.CommentLike{}).
		Select("comment_id, count(*) as count").
		Where("comment_id IN ?", commentIDs).
		Group("comment_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
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
		Where("user_id = ? AND comment_id IN ?", userID, commentIDs).
		Find(&likes).Error; err != nil {
		return nil, err
	}

	for _, l := range likes {
		res[l.CommentID] = true
	}
	return res, nil
}
