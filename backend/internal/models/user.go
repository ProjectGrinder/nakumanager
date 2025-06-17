package model

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string
	Email        string `json:"email"`
}
