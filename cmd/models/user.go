package models

type User struct {
	ID       int
	Email    string
	Password string
}

// UserLoginRequest represents the request payload for logging in a user
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required" example:"test@gmail.com"`
	Password string `json:"password" binding:"required" example:"password"`
}

// PasswordResetRequest represents the request payload for resetting a password
type PasswordResetRequest struct {
	CurrentPass string `json:"current_password" binding:"required" example:"oldpassword"`
	NewPass     string `json:"new_password" binding:"required,min=8" example:"newsecurepassword"`
}
