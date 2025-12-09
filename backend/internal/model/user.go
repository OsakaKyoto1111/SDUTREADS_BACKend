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

	PostsCount     int `gorm:"column:posts_count;default:0"`
	FollowersCount int `gorm:"column:followers_count;default:0"`
	FollowingCount int `gorm:"column:following_count;default:0"`

	Posts     []Post     `gorm:"foreignKey:UserID"`
	PostLikes []PostLike `gorm:"foreignKey:UserID"`
	Comments  []Comment  `gorm:"foreignKey:UserID"`
	Followers []Follower `gorm:"foreignKey:UserID"`
	Following []Follower `gorm:"foreignKey:FollowerID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
