package repository

import (
	"fmt"
	"strings"

	"backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNotFound = fmt.Errorf("not found")

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByIDWithPreloads(id uint) (*model.User, error)
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
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByIDWithPreloads(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Posts").Preload("Followers").Preload("Following").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user with preloads: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get by email: %w", err)
	}
	return &u, nil
}

func (r *userRepository) GetByNickname(nickname string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("nickname = ?", nickname).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get by nickname: %w", err)
	}
	return &u, nil
}

func (r *userRepository) Update(user *model.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (r *userRepository) GetPostsCount(userID uint) (int64, error) {
	var cnt int64
	if err := r.db.Model(&model.Post{}).Where("user_id = ?", userID).Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("count posts: %w", err)
	}
	return cnt, nil
}

func (r *userRepository) GetFollowersCount(userID uint) (int64, error) {
	var cnt int64
	if err := r.db.Model(&model.Follower{}).Where("user_id = ?", userID).Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("count followers: %w", err)
	}
	return cnt, nil
}

func (r *userRepository) GetFollowingCount(userID uint) (int64, error) {
	var cnt int64
	if err := r.db.Model(&model.Follower{}).Where("follower_id = ?", userID).Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("count following: %w", err)
	}
	return cnt, nil
}

func (r *userRepository) Search(query string) ([]model.User, error) {
	var users []model.User
	query = strings.TrimSpace(query)
	if query == "" {
		return users, nil
	}

	words := strings.Fields(query) // Разбиваем на слова
	dbQuery := r.db.Model(&model.User{})

	// Основной WHERE
	for _, word := range words {
		wordPattern := "%" + word + "%"
		dbQuery = dbQuery.Where(
			"nickname ILIKE ? OR COALESCE(first_name, '') ILIKE ? OR COALESCE(last_name, '') ILIKE ? OR "+
				"COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') ILIKE ? OR "+
				"COALESCE(last_name, '') || ' ' || COALESCE(first_name, '') ILIKE ?",
			wordPattern, wordPattern, wordPattern, wordPattern, wordPattern,
		)
	}

	// Релевантность через встроенную строку
	fullName := strings.Join(words, " ")
	fullNameReversed := strings.Join(reverseSlice(words), " ")
	relevanceOrder := fmt.Sprintf(`
		CASE
			WHEN nickname ILIKE '%s' THEN 1
			WHEN COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') ILIKE '%s' THEN 2
			ELSE 3
		END
	`, fullName, fullNameReversed)

	dbQuery = dbQuery.Order(relevanceOrder)

	if err := dbQuery.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("search users: %w", err)
	}

	return users, nil
}

func reverseSlice(s []string) []string {
	reversed := make([]string, len(s))
	for i, word := range s {
		reversed[len(s)-1-i] = word
	}
	return reversed
}

func (r *userRepository) Follow(userID, targetID uint) error {
	if userID == 0 || targetID == 0 {
		return fmt.Errorf("invalid ids")
	}
	if userID == targetID {
		return fmt.Errorf("cannot follow yourself")
	}
	f := model.Follower{UserID: targetID, FollowerID: userID}
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&f).Error; err != nil {
		return fmt.Errorf("follow create: %w", err)
	}
	return nil
}

func (r *userRepository) Unfollow(userID, targetID uint) error {
	if err := r.db.Where("user_id = ? AND follower_id = ?", targetID, userID).
		Delete(&model.Follower{}).Error; err != nil {
		return fmt.Errorf("unfollow delete: %w", err)
	}
	return nil
}
