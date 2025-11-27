package repository

import (
	"context"
	"errors"

	"backend/internal/model"

	"gorm.io/gorm"
)

var (
	// ErrNotFound signals any repository query that cannot find a record.
	ErrNotFound = errors.New("record not found")
)

// UserRepository defines operations against the user storage.
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByNickname(ctx context.Context, nickname string) (*model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	UpdateProfile(ctx context.Context, id uint, fields map[string]interface{}) error
	SearchUsers(ctx context.Context, query string, limit int) ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a repository tied to the provided DB connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &user, err
}

func (r *userRepository) FindByNickname(ctx context.Context, nickname string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("nickname = ?", nickname).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &user, err
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &user, err
}

func (r *userRepository) UpdateProfile(ctx context.Context, id uint, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *userRepository) SearchUsers(ctx context.Context, query string, limit int) ([]model.User, error) {
	if limit <= 0 {
		limit = 25
	}
	pattern := "%" + query + "%"
	var users []model.User
	err := r.db.WithContext(ctx).
		Where("nickname ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?", pattern, pattern, pattern).
		Limit(limit).
		Find(&users).
		Error
	return users, err
}
