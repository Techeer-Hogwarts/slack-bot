package models

type User struct {
	ID       int
	Email    string
	Password string
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
