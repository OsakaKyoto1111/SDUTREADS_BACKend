package model

type File struct {
	ID     uint   `gorm:"primaryKey"`
	PostID uint   `gorm:"index;not null"`
	Post   Post   `gorm:"foreignKey:PostID"`
	URL    string `gorm:"not null"`
}
