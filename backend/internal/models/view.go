package model

type CreateView struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	TeamID    string   `json:"team_id"`
	Assignnee string   `json:"assignee"`
	GroupBys  []string `json:"group_bys"`
}

type ViewGroupBy struct {
	ViewID string `json:"view_id"`
}

type EditView struct {
	ID       string   `json:"view_id" validate:"required"`
	TeamID   string   `json:"team_id"`
	Name     string   `json:"name"`
	GroupBys []string `json:"group_bys"`
}
