package model

type CreateView struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	TeamID    string   `json:"team_id"`
	Assignnee string   `json:"assignee"`
	GroupBys  []string `json:"group_bys"`
}

type ViewGroupBy struct {
	GroupBys []string `json:"group_bys"`
}

type UpdateViewRequest struct {
	TeamID   string   `json:"team_id"`
	Name     string   `json:"name"`
	GroupBys []string `json:"group_bys"`
}
