package mapper

import (
	"backend/internal/dto"
	"backend/internal/model"
	"time"
)

func MapPostsToDTO(posts []model.Post, userID uint) []dto.PostResponse {
	result := make([]dto.PostResponse, 0, len(posts))

	for _, p := range posts {
		files := make([]dto.FileResponse, 0, len(p.Files))
		for _, f := range p.Files {
			files = append(files, dto.FileResponse{
				ID:  f.ID,
				URL: f.URL,
			})
		}

		isLiked := false
		for _, l := range p.Likes {
			if l.UserID == userID {
				isLiked = true
				break
			}
		}

		author := dto.PostAuthorDTO{
			ID:        p.User.ID,
			Nickname:  p.User.Nickname,
			AvatarURL: p.User.AvatarURL,
		}

		result = append(result, dto.PostResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			User:        author,
			Description: p.Description,
			Files:       files,
			LikesCount:  len(p.Likes),
			Comments:    len(p.Comments),
			IsLiked:     isLiked,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		})
	}

	return result
}
