package service

import (
	"backend/internal/dto"
	"backend/internal/mapper"
	"backend/internal/model"
	"backend/internal/repository"
	"errors"
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
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) SearchUsersWithCounts(query string) ([]dto.UserResponse, error) {
	users, err := s.repo.Search(query)
	if err != nil {
		return nil, err
	}

	var usersResp []dto.UserResponse
	for _, u := range users {
		postsCnt, _ := s.repo.GetPostsCount(u.ID)
		followersCnt, _ := s.repo.GetFollowersCount(u.ID)
		followingCnt, _ := s.repo.GetFollowingCount(u.ID)
		usersResp = append(usersResp, mapper.MapUserToResponseWithCounts(&u, postsCnt, followersCnt, followingCnt))
	}

	return usersResp, nil
}

func (s *userService) GetUser(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	postsCnt, err := s.repo.GetPostsCount(id)
	if err != nil {
		return nil, err
	}
	followersCnt, err := s.repo.GetFollowersCount(id)
	if err != nil {
		return nil, err
	}
	followingCnt, err := s.repo.GetFollowingCount(id)
	if err != nil {
		return nil, err
	}

	resp := mapper.MapUserToResponseWithCounts(user, postsCnt, followersCnt, followingCnt)
	return &resp, nil
}

func (s *userService) UpdateUser(id uint, dto dto.UpdateUserDTO) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	mapper.ApplyUpdateUserDTO(user, dto)

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

func (s *userService) SearchUsers(query string) ([]model.User, error) {
	return s.repo.Search(query)
}

func (s *userService) SetAvatar(id uint, avatarURL string) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.AvatarURL = &avatarURL
	err = s.repo.Update(user)
	return user, err
}

func (s *userService) Follow(userID uint, targetID uint) error {
	if userID == targetID {
		return errors.New("cannot follow yourself")
	}
	return s.repo.Follow(userID, targetID)
}

func (s *userService) Unfollow(userID uint, targetID uint) error {
	return s.repo.Unfollow(userID, targetID)
}
