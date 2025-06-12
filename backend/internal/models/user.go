package model

type User struct {
	Username     string `json:"username"`
	PasswordHash string
	Email        string `json:"email"`
}
