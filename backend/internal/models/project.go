package model

type CreateProjectRequest struct {
	Name        string  `json:"name"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	WorkspaceID string  `json:"workspace_id"`
	TeamID      string  `json:"team_id"`
	LeaderID    *string `json:"leader_id"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Label       *string `json:"label"`
	CreatedBy   string  `json:"created_by"`
}
