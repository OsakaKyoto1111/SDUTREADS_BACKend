package dto

import "time"

// UserResponse is a safe view of the user entity without credentials.
type UserResponse struct {
	ID               uint      `json:"id"`
	Email            string    `json:"email"`
	Nickname         string    `json:"nickname"`
	FirstName        *string   `json:"first_name,omitempty"`
	LastName         *string   `json:"last_name,omitempty"`
	AvatarURL        *string   `json:"avatar_url,omitempty"`
	Grade            *string   `json:"grade,omitempty"`
	Major            *string   `json:"major,omitempty"`
	City             *string   `json:"city,omitempty"`
	Description      *string   `json:"description,omitempty"`
	PostsCount       int       `json:"posts_count"`
	SubscribersCount int       `json:"subscribers_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// UpdateProfileRequest carries optional profile attributes for patching.
type UpdateProfileRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Nickname    *string `json:"nickname,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Grade       *string `json:"grade,omitempty"`
	Major       *string `json:"major,omitempty"`
	City        *string `json:"city,omitempty"`
	Description *string `json:"description,omitempty"`
}
