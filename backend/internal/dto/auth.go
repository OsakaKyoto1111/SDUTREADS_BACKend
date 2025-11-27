package dto

// RegisterRequest describes a payload for new user registration.
type RegisterRequest struct {
	Email     string  `json:"email"`
	Nickname  string  `json:"nickname"`
	Password  string  `json:"password"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Grade     *string `json:"grade,omitempty"`
	Major     *string `json:"major,omitempty"`
	City      *string `json:"city,omitempty"`
}

// LoginRequest describes credentials for login.
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

// AuthResponse combines user view and access token.
type AuthResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}
