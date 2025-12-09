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

type CommentTree struct {
	ID        uint          `json:"id"`
	PostID    uint          `json:"post_id"`
	UserID    uint          `json:"user_id"`
	ParentID  *uint         `json:"parent_id"`
	Text      string        `json:"text"`
	Likes     int           `json:"likes"`
	IsLiked   bool          `json:"is_liked"`
	CreatedAt time.Time     `json:"created_at"`
	User      UserShortDTO  `json:"user"`
	Replies   []CommentTree `json:"replies"`
}

type PostWithCommentsResponse struct {
	ID          uint           `json:"id"`
	Description *string        `json:"description"`
	Files       []FileResponse `json:"files"`
	LikesCount  int            `json:"likes_count"`
	IsLiked     bool           `json:"is_liked"`
	Comments    []CommentTree  `json:"comments"`
}
