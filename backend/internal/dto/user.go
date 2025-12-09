package dto

type UpdateUserDTO struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	AvatarURL   *string `json:"avatar_url"`
	Grade       *string `json:"grade"`
	Major       *string `json:"major"`
	City        *string `json:"city"`
	Description *string `json:"description"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`

	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	AvatarURL   *string `json:"avatar_url"`
	Grade       *string `json:"grade"`
	Major       *string `json:"major"`
	City        *string `json:"city"`
	Description *string `json:"description"`

	PostsCount     int `json:"posts_count"`
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
}
