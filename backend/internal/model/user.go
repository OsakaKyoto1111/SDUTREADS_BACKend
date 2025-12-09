package model

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	Nickname     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`

	FirstName   *string
	LastName    *string
	AvatarURL   *string
	Grade       *string
	Major       *string
	City        *string
	Description *string

	// followers: все записи followers, где followers.user_id == user.ID
	Followers []Follower `gorm:"foreignKey:UserID"`

	// following: все записи followers, где followers.follower_id == user.ID
	Following []Follower `gorm:"foreignKey:FollowerID"`

	Posts     []Post     `gorm:"foreignKey:UserID"`
	PostLikes []PostLike `gorm:"foreignKey:UserID"`
	Comments  []Comment  `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
