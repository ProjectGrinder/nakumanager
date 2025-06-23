package model

type Workspace struct {
	ID      string   `json:"id"`
	Name    string   `json:"name" validate:"required"`
	Members []string `json:"members"`
}
