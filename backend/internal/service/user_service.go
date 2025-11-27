package service

import (
	"context"
	"errors"
	"strings"

	"backend/internal/dto"
	"backend/internal/repository"
)

// UserService manages user profile operations.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService constructs a UserService.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetMe returns the current user.
func (s *UserService) GetMe(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := mapToUserResponse(user)
	return &resp, nil
}

// UpdateProfile updates mutable profile fields and returns the fresh view.
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	updates := make(map[string]interface{})

	if req.Nickname != nil {
		nickname := strings.TrimSpace(*req.Nickname)
		if nickname == "" {
			return nil, errors.New("nickname must not be empty")
		}
		existing, err := s.repo.FindByNickname(ctx, nickname)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
		if err == nil && existing.ID != userID {
			return nil, errors.New("nickname already taken")
		}
		updates["nickname"] = nickname
	}

	if req.FirstName != nil {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = req.LastName
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Grade != nil {
		updates["grade"] = req.Grade
	}
	if req.Major != nil {
		updates["major"] = req.Major
	}
	if req.City != nil {
		updates["city"] = req.City
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}

	if err := s.repo.UpdateProfile(ctx, userID, updates); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := mapToUserResponse(user)
	return &resp, nil
}

// SearchUsers finds users matching nickname, first name or last name.
func (s *UserService) SearchUsers(ctx context.Context, query string, limit int) ([]dto.UserResponse, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("query is required")
	}
	users, err := s.repo.SearchUsers(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	results := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		results = append(results, mapToUserResponse(&u))
	}
	return results, nil
}

// GetUserByID returns another user's profile.
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := mapToUserResponse(user)
	return &resp, nil
}
