package dto

type FeedPostItem struct {
	ID          uint           `json:"id"`
	User        UserShortDTO   `json:"user"`
	Description *string        `json:"description,omitempty"`
	Files       []FileResponse `json:"files"`
	LikesCount  int            `json:"likes_count"`
	Comments    int            `json:"comments"`
	IsLiked     bool           `json:"is_liked"`
	CreatedAt   string         `json:"created_at"`
}

type FeedResponse struct {
	Posts      []PostResponse `json:"posts"`
	NextCursor *string        `json:"next_cursor,omitempty"`
	HasMore    bool           `json:"has_more"`
}
