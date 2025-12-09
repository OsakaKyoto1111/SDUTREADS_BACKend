package model

type Follower struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"index;not null"`
	User   User `gorm:"foreignKey:UserID"`

	FollowerID uint `gorm:"index;not null"`
	Follower   User `gorm:"foreignKey:FollowerID"`
}
