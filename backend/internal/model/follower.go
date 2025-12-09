package model

type Follower struct {
	ID uint `gorm:"primaryKey"`

	// user_id — профиль, на который подписываются (target)
	UserID uint  `gorm:"index;not null;column:user_id"`
	User   *User `gorm:"foreignKey:UserID;references:ID" json:"-"`

	// follower_id — кто подписался
	FollowerID uint  `gorm:"index;not null;column:follower_id"`
	Follower   *User `gorm:"foreignKey:FollowerID;references:ID" json:"-"`
}
