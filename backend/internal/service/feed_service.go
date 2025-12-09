package service

import (
	"backend/internal/dto"
	"backend/internal/mapper"
	"backend/internal/repository"
	"time"
)

type FeedService struct {
	repo *repository.FeedRepository
}

func NewFeedService(repo *repository.FeedRepository) *FeedService {
	return &FeedService{repo: repo}
}

const MixRatio = 4

func (s *FeedService) GetFeed(userID uint, limit int, cursor *time.Time) (*dto.FeedResponse, error) {
	// Following posts
	following, err := s.repo.GetFollowingPosts(userID, limit, cursor)
	if err != nil {
		return nil, err
	}

	// Recommended
	recCount := limit / MixRatio
	if recCount < 1 {
		recCount = 1
	}

	recommended, _ := s.repo.GetRecommendedPosts(userID, recCount)

	// Mapping
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

	// cursor logic
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
