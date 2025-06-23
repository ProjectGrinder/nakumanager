package model

type Team struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Members []User `json:"user"`
	Leader  User   `json:"leader"`
}
