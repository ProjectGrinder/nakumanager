package model

type View struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	UserID   uint     `json:"user_id"`
	TeamID   uint     `json:"team_id"`
	GroupBys []string `json:"group_bys"`
	Issues   []Issue  `json:"issues"`
}
