package dto

import "time"

type AddCommentRequest struct {
	Text     string `json:"text" validate:"required"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

type CommentDTO struct {
	ID        uint         `json:"id"`
	PostID    uint         `json:"post_id"`
	UserID    uint         `json:"user_id"`
	ParentID  *uint        `json:"parent_id,omitempty"`
	Text      string       `json:"text"`
	Likes     int          `json:"likes"`
	IsLiked   bool         `json:"is_liked"`
	User      UserShortDTO `json:"user"`
	Replies   []CommentDTO `json:"replies,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

type UserShortDTO struct {
	ID       uint    `json:"id"`
	Nickname string  `json:"nickname"`
	Avatar   *string `json:"avatar,omitempty"`
}

type CommentListResponse struct {
	Comments []CommentDTO `json:"comments"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	Total    int64        `json:"total"`
}
