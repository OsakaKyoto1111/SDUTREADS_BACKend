package dto

type CreatePostRequest struct {
	Description *string `json:"description"`
}

type UpdatePostRequest struct {
	Description *string `json:"description"`
}

type PostResponse struct {
	ID          uint           `json:"id"`
	Description *string        `json:"description"`
	Files       []FileResponse `json:"files"`
	LikesCount  int            `json:"likes_count"`
	Comments    []CommentRes   `json:"comments"`
}

type FileResponse struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

type CommentRes struct {
	ID        uint         `json:"id"`
	UserID    uint         `json:"user_id"`
	PostID    uint         `json:"post_id"`
	ParentID  *uint        `json:"parent_id"`
	Text      string       `json:"text"`
	Likes     int          `json:"likes"`
	Replies   []CommentRes `json:"replies"`
	CreatedAt string       `json:"created_at"`
}
