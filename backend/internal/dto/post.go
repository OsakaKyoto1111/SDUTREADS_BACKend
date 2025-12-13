package dto

type CreatePostRequest struct {
	Description *string `json:"description"`
}

type UpdatePostRequest struct {
	Description *string `json:"description"`
}

type PostResponse struct {
	ID          uint           `json:"id"`
	UserID      uint           `json:"user_id"`
	User        PostAuthorDTO  `json:"user"`
	Description *string        `json:"description,omitempty"`
	Files       []FileResponse `json:"files"`
	LikesCount  int            `json:"likes_count"`
	Comments    int            `json:"comments"`
	IsLiked     bool           `json:"is_liked"`
	CreatedAt   string         `json:"created_at"`
}

type FileResponse struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

type CreatePostRequestMultipart struct {
	Description *string `form:"description"`
}

// PostAuthorDTO contains only the fields the feed needs to render author info.
type PostAuthorDTO struct {
	ID        uint    `json:"id"`
	Nickname  string  `json:"nickname"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
