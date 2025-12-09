package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Nickname string `json:"nickname" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`

	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Grade       *string `json:"grade"`
	Major       *string `json:"major"`
	City        *string `json:"city"`
	Description *string `json:"description"`
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}
