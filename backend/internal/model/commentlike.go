package model

type CommentLike struct {
	ID        uint    `gorm:"primaryKey"`
	CommentID uint    `gorm:"index;not null"`
	Comment   Comment `gorm:"foreignKey:CommentID"`

	UserID uint `gorm:"index;not null"`
	User   User `gorm:"foreignKey:UserID"`
}
