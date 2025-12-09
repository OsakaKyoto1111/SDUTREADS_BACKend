package mapper

import (
	"backend/internal/dto"
	"backend/internal/model"
)

func MapRegisterRequestToUser(req dto.RegisterRequest, passwordHash string) *model.User {
	return &model.User{
		Email:        req.Email,
		Nickname:     req.Nickname,
		PasswordHash: passwordHash,

		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Grade:       req.Grade,
		Major:       req.Major,
		City:        req.City,
		Description: req.Description,
	}
}

func MapUserToResponseWithCounts(u *model.User, postsCount, followersCount, followingCount int64) dto.UserResponse {
	return dto.UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Nickname: u.Nickname,

		FirstName:   u.FirstName,
		LastName:    u.LastName,
		AvatarURL:   u.AvatarURL,
		Grade:       u.Grade,
		Major:       u.Major,
		City:        u.City,
		Description: u.Description,

		PostsCount:     int(postsCount),
		FollowersCount: int(followersCount),
		FollowingCount: int(followingCount),
	}
}

func ApplyUpdateUserDTO(u *model.User, dto dto.UpdateUserDTO) {
	if dto.FirstName != nil {
		u.FirstName = dto.FirstName
	}
	if dto.LastName != nil {
		u.LastName = dto.LastName
	}
	if dto.AvatarURL != nil {
		u.AvatarURL = dto.AvatarURL
	}
	if dto.Grade != nil {
		u.Grade = dto.Grade
	}
	if dto.Major != nil {
		u.Major = dto.Major
	}
	if dto.City != nil {
		u.City = dto.City
	}
	if dto.Description != nil {
		u.Description = dto.Description
	}
}
