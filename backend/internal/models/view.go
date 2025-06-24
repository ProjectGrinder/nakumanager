package model

type ViewCreateRequest struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	TeamID    string   `json:"team_id"`
	Assignnee string   `json:"assignee"`
	GroupBys  []string `json:"group_bys"`
}

type ViewGroupBy struct {
	ViewID string `json:"view_id"`
}
