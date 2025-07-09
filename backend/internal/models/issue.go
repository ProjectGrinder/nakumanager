package model

import (
	"time"
)

type IssueCreate struct {
	ID        string     `json:"id"`
	Title     string     `json:"title" validate:"required"`
	Content   *string    `json:"content,omitempty"`
	Priority  *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	Status    string     `json:"status,omitempty" validate:"omitempty,oneof=todo doing done"`
	Assignee  *[]string  `json:"assignee,omitempty"`
	ProjectID *string    `json:"project_id,omitempty"`
	TeamID    string     `json:"team_id" validate:"required"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Label     *string    `json:"label,omitempty"`
	OwnerID   string     `json:"owner_id" validate:"required"`
}

type UpdateIssueRequest struct {
	ID             string     `json:"issue_id"`
	Title          *string    `json:"title,omitempty"`
	Content        *string    `json:"content,omitempty"`
	Priority       *string    `json:"priority,omitempty"`
	Status         *string    `json:"status,omitempty"`
	AddAssignee    *[]string  `json:"add_assignees,omitempty"`
	RemoveAssignee *[]string  `json:"remove_assignees,omitempty"`
	ProjectID      *string    `json:"project_id,omitempty"`
	TeamID         *string    `json:"team_id,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Label          *string    `json:"label,omitempty"`
	OwnerID        *string    `json:"owner_id,omitempty"`
}
