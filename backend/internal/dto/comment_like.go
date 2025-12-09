package dto

type CommentLikeResponse struct {
	CommentID uint `json:"comment_id"`
	Liked     bool `json:"liked"`
}
