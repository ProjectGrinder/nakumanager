package model

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string
	Email        string `json:"email"`
}
