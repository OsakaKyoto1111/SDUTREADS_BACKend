package service

import (
	"backend/internal/dto"
	"backend/internal/model"
)

func mapToUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		Nickname:       user.Nickname,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AvatarURL:      user.AvatarURL,
		Grade:          user.Grade,
		Major:          user.Major,
		City:           user.City,
		Description:    user.Description,
		PostsCount:     user.PostsCount,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
