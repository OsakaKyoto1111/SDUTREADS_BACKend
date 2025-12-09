package model

import "time"

type Comment struct {
	ID uint `gorm:"primaryKey"`

	PostID uint `gorm:"index;not null"`
	Post   Post `gorm:"foreignKey:PostID"`

	UserID uint `gorm:"index;not null"`
	User   User `gorm:"foreignKey:UserID"`

	ParentID *uint
	Parent   *Comment `gorm:"foreignKey:ParentID"`

	Text string `gorm:"type:text;not null"`

	Replies []Comment `gorm:"foreignKey:ParentID"`

	Likes []CommentLike `gorm:"foreignKey:CommentID"`

	CreatedAt time.Time
}
