package dto

type CreatePostRequest struct {
	Description *string `json:"description"`
}

type UpdatePostRequest struct {
	Description *string `json:"description"`
}

type PostResponse struct {
	ID          uint           `json:"id"`
	Description *string        `json:"description,omitempty"`
	Files       []FileResponse `json:"files"`
	LikesCount  int            `json:"likes_count"`
	Comments    int            `json:"comments"`
	IsLiked     bool           `json:"is_liked"`
}

type FileResponse struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

type CreatePostRequestMultipart struct {
	Description *string `form:"description"`
}
