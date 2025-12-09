package model

type PostLike struct {
	ID     uint `gorm:"primaryKey"`
	PostID uint `gorm:"index;not null"`
	Post   Post `gorm:"foreignKey:PostID"`

	UserID uint `gorm:"index;not null"`
	User   User `gorm:"foreignKey:UserID"`
}
