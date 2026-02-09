package dto

// RegisterRequest is the DTO for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest is the DTO for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is the DTO for auth responses (login/register)
type AuthResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

// UserProfile is the public user data (no password)
type UserProfile struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsOnline bool   `json:"is_online"`
}
