package model

import "time"

type Post struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"index;not null"`
	User   User `gorm:"foreignKey:UserID"`

	Description *string
	ViewsCount  int `gorm:"default:0"`

	Files    []File     `gorm:"foreignKey:PostID"`
	Likes    []PostLike `gorm:"foreignKey:PostID"`
	Comments []Comment  `gorm:"foreignKey:PostID"`

	CreatedAt time.Time
}
