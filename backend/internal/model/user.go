package model

import "time"

// User represents the persisted user entity.
type User struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Email            string    `gorm:"uniqueIndex;not null" json:"email"`
	Nickname         string    `gorm:"uniqueIndex;not null" json:"nickname"`
	PasswordHash     string    `gorm:"not null" json:"-"`
	FirstName        *string   `gorm:"type:text" json:"first_name,omitempty"`
	LastName         *string   `gorm:"type:text" json:"last_name,omitempty"`
	AvatarURL        *string   `gorm:"type:text" json:"avatar_url,omitempty"`
	Grade            *string   `gorm:"type:text" json:"grade,omitempty"`
	Major            *string   `gorm:"type:text" json:"major,omitempty"`
	City             *string   `gorm:"type:text" json:"city,omitempty"`
	Description      *string   `gorm:"type:text" json:"description,omitempty"`
	PostsCount       int       `gorm:"default:0" json:"posts_count"`
	SubscribersCount int       `gorm:"default:0" json:"subscribers_count"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
