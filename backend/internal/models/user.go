package model

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username, validate:"required"`
	PasswordHash string `validate:"required"`
	Email        string `json:"email, validate:"required"`
	Roles        string `json:"roles, validate:"required"`
}

func GetUserByID(id string) User {
	//TODO: ดึงจาก database
	return User{ID: id}
}
