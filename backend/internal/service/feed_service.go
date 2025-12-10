package service

import (
	"backend/internal/model"
	"fmt"
	"time"

	"backend/internal/dto"
	"backend/internal/mapper"
	"backend/internal/repository"
)

type FeedService interface {
	GetFeed(userID uint, limit int, cursor *time.Time) (*dto.FeedResponse, error)
}

type feedService struct {
	repo repository.FeedRepository
}

func NewFeedService(r repository.FeedRepository) FeedService {
	return &feedService{repo: r}
}

const MixRatio = 4

func (s *feedService) GetFeed(userID uint, limit int, cursor *time.Time) (*dto.FeedResponse, error) {
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}
	if limit <= 0 {
		limit = 20
	}

	following, err := s.repo.GetFollowingPosts(userID, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("get following posts: %w", err)
	}

	recCount := limit / MixRatio
	if recCount < 1 {
		recCount = 1
	}

	recommended, err := s.repo.GetRecommendedPosts(userID, recCount)
	if err != nil {
		recommended = []model.Post{}
	}

	followingDTO := mapper.MapPostsToDTO(following, userID)
	recommendedDTO := mapper.MapPostsToDTO(recommended, userID)

	result := []dto.PostResponse{}
	recIndex := 0

	for i, p := range followingDTO {
		result = append(result, p)

		if (i+1)%MixRatio == 0 && recIndex < len(recommendedDTO) {
			result = append(result, recommendedDTO[recIndex])
			recIndex++
		}
	}

	for recIndex < len(recommendedDTO) {
		result = append(result, recommendedDTO[recIndex])
		recIndex++
	}

	var nextCursor *string
	if len(following) > 0 {
		t := following[len(following)-1].CreatedAt.UTC().Format(time.RFC3339)
		nextCursor = &t
	}

	return &dto.FeedResponse{
		Posts:      result,
		NextCursor: nextCursor,
		HasMore:    len(following) == limit,
	}, nil
}
