package model

type Team struct {
	ID     string `json:"id" validate:"required,uuid"`
	Name   string `json:"name" validate:"required"`
	Leader User   `json:"leader"`
}

type EditTeam struct {
	TeamID       string `json:"team_id"`
	Name         string `json:"name"`
	AddMember    string `json:"add_member"`
	RemoveMember string `json:"remove_member"`
	Leader       string `json:"leader"`
}
