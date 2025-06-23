package model

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username" validate:"required"`
	PasswordHash string `json:"password_hash" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Roles        string `json:"roles" validate:"required"`
}

