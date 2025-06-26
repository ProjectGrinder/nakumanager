package model

import "time"

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

type EditProject struct {
	ID           string     `json:"project_id"`
	Name         string     `json:"name"`
	LeaderID     string     `json:"leader_id"`
	Status       *string    `json:"status"`
	Priority     *string    `json:"priority"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Label        *string    `json:"label"`
	AddMember    string     `json:"add_member"`
	RemoveMember string     `json:"remove_member"`
	WorkspaceID  string     `json:"workspace_id"`
}
