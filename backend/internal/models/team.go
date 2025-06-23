package model

type Team struct {
	ID     string `json:"id" validate:"required,uuid"`
	Name   string `json:"name" validate:"required"`
	Leader User   `json:"leader"`
}
