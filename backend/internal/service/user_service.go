package service

import (
	"fmt"

	"backend/internal/dto"
	"backend/internal/mapper"
	"backend/internal/model"
	"backend/internal/repository"
)

type UserService interface {
	GetUser(id uint) (*dto.UserResponse, error)
	UpdateUser(id uint, dto dto.UpdateUserDTO) (*model.User, error)
	DeleteUser(id uint) error
	SearchUsersWithCounts(query string) ([]dto.UserResponse, error)
	SearchUsers(query string) ([]model.User, error)
	SetAvatar(id uint, avatarURL string) (*model.User, error)
	Follow(userID uint, targetID uint) error
	Unfollow(userID uint, targetID uint) error
	GetFollowers(userID uint) ([]dto.UserResponse, error)
	GetFollowing(userID uint) ([]dto.UserResponse, error)
	IsFollowing(userID uint, targetID uint) (bool, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
func (s *userService) IsFollowing(userID uint, targetID uint) (bool, error) {
	if userID == 0 || targetID == 0 {
		return false, fmt.Errorf("invalid ids")
	}
	if userID == targetID {
		return false, nil
	}

	return s.repo.IsFollowing(userID, targetID)
}

func (s *userService) SearchUsersWithCounts(query string) ([]dto.UserResponse, error) {
	users, err := s.repo.Search(query)
	if err != nil {
		return nil, fmt.Errorf("search users: %w", err)
	}

	var usersResp []dto.UserResponse
	for _, u := range users {
		postsCnt, _ := s.repo.GetPostsCount(u.ID)
		followersCnt, _ := s.repo.GetFollowersCount(u.ID)
		followingCnt, _ := s.repo.GetFollowingCount(u.ID)
		usersResp = append(usersResp,
			mapper.MapUserToResponseWithCounts(&u, postsCnt, followersCnt, followingCnt),
		)
	}

	return usersResp, nil
}

func (s *userService) GetUser(id uint) (*dto.UserResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid id")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	postsCnt, _ := s.repo.GetPostsCount(id)
	followersCnt, _ := s.repo.GetFollowersCount(id)
	followingCnt, _ := s.repo.GetFollowingCount(id)

	resp := mapper.MapUserToResponseWithCounts(user, postsCnt, followersCnt, followingCnt)
	return &resp, nil
}

func (s *userService) UpdateUser(id uint, dto dto.UpdateUserDTO) (*model.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid id")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	mapper.ApplyUpdateUserDTO(user, dto)

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid id")
	}
	return s.repo.Delete(id)
}

func (s *userService) SearchUsers(query string) ([]model.User, error) {
	return s.repo.Search(query)
}

func (s *userService) SetAvatar(id uint, avatarURL string) (*model.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid id")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	user.AvatarURL = &avatarURL

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("set avatar: %w", err)
	}

	return user, nil
}

func (s *userService) Follow(userID uint, targetID uint) error {
	if userID == 0 || targetID == 0 {
		return fmt.Errorf("invalid ids")
	}
	if userID == targetID {
		return fmt.Errorf("cannot follow yourself")
	}
	return s.repo.Follow(userID, targetID)
}

func (s *userService) Unfollow(userID uint, targetID uint) error {
	if userID == 0 || targetID == 0 {
		return fmt.Errorf("invalid ids")
	}
	return s.repo.Unfollow(userID, targetID)
}

func (s *userService) GetFollowers(userID uint) ([]dto.UserResponse, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid id")
	}

	users, err := s.repo.GetFollowers(userID)
	if err != nil {
		return nil, fmt.Errorf("get followers: %w", err)
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		postsCnt, _ := s.repo.GetPostsCount(u.ID)
		followersCnt, _ := s.repo.GetFollowersCount(u.ID)
		followingCnt, _ := s.repo.GetFollowingCount(u.ID)

		resp = append(resp, mapper.MapUserToResponseWithCounts(&u, postsCnt, followersCnt, followingCnt))
	}

	return resp, nil
}

func (s *userService) GetFollowing(userID uint) ([]dto.UserResponse, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid id")
	}

	users, err := s.repo.GetFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("get following: %w", err)
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		postsCnt, _ := s.repo.GetPostsCount(u.ID)
		followersCnt, _ := s.repo.GetFollowersCount(u.ID)
		followingCnt, _ := s.repo.GetFollowingCount(u.ID)

		resp = append(resp, mapper.MapUserToResponseWithCounts(&u, postsCnt, followersCnt, followingCnt))
	}

	return resp, nil
}
