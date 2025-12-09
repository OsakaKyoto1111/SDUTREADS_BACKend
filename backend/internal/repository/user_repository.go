package repository

import (
	"backend/internal/model"
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByIDWithPreloads(id uint) (*model.User, error) // optional
	GetByEmail(email string) (*model.User, error)
	GetByNickname(nickname string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error

	GetPostsCount(userID uint) (int64, error)
	GetFollowersCount(userID uint) (int64, error)
	GetFollowingCount(userID uint) (int64, error)

	Search(query string) ([]model.User, error)
	Follow(userID, targetID uint) error
	Unfollow(userID, targetID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByIDWithPreloads(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Posts").Preload("Followers").Preload("Following").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetPostsCount(userID uint) (int64, error) {
	var cnt int64
	err := r.db.Model(&model.Post{}).Where("user_id = ?", userID).Count(&cnt).Error
	return cnt, err
}

func (r *userRepository) GetFollowersCount(userID uint) (int64, error) {
	var cnt int64
	err := r.db.Model(&model.Follower{}).Where("user_id = ?", userID).Count(&cnt).Error
	return cnt, err
}

func (r *userRepository) GetFollowingCount(userID uint) (int64, error) {
	var cnt int64
	err := r.db.Model(&model.Follower{}).Where("follower_id = ?", userID).Count(&cnt).Error
	return cnt, err
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *userRepository) GetByNickname(nickname string) (*model.User, error) {
	var u model.User
	err := r.db.Where("nickname = ?", nickname).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) Search(query string) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("nickname ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&users).Error
	return users, err
}

func (r *userRepository) Follow(userID, targetID uint) error {
	var existing model.Follower
	err := r.db.Where("user_id = ? AND follower_id = ?", targetID, userID).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return r.db.Create(&model.Follower{
		UserID:     targetID,
		FollowerID: userID,
	}).Error
}

func (r *userRepository) Unfollow(userID, targetID uint) error {
	return r.db.Where("user_id = ? AND follower_id = ?", targetID, userID).Delete(&model.Follower{}).Error
}
