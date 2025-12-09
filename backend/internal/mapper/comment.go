package mapper

import (
	"backend/internal/dto"
	"backend/internal/model"
)

func MapCommentToDTO(c model.Comment, isLiked bool, likesCount int) dto.CommentDTO {
	// Replies recursive mapping
	var replies []dto.CommentDTO
	for _, r := range c.Replies {
		replies = append(replies, MapCommentToDTO(r, false, len(r.Likes)))
	}

	return dto.CommentDTO{
		ID:       c.ID,
		PostID:   c.PostID,
		UserID:   c.UserID,
		ParentID: c.ParentID,
		Text:     c.Text,
		Likes:    likesCount,
		IsLiked:  isLiked,
		User: dto.UserShortDTO{
			ID:       c.User.ID,
			Nickname: c.User.Nickname,
			Avatar:   c.User.AvatarURL,
		},
		Replies:   replies,
		CreatedAt: c.CreatedAt,
	}
}
