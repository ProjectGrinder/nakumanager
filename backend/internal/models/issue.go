package model

import (
	"time"
)

type IssueCreate struct {
	ID        string     `json:"id"`
	Title     string     `json:"title" validate:"required"`
	Content   *string    `json:"content,omitempty"`
	Priority  *string    `json:"priority" validate:"omitempty,oneof=low medium high"`
	Status    string     `json:"status" validate:"omitempty,oneof=todo doing done"`
	Assignee  *string    `json:"assignee,omitempty"`
	ProjectID *string    `json:"project_id" validate:"required"`
	TeamID    string     `json:"team_id" validate:"required"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Label     *string    `json:"label,omitempty"`
	OwnerID   string     `json:"ownerId" validate:"required"`
}

type AssigneeRequest struct {
	IssueID string `json:"issue_id"`
	UserID  string `json:"user_id"`
}

type EditIssue struct {
	ID        string     `json:"issue_id"`
	Title     *string    `json:"title,omitempty"`
	Content   *string    `json:"content,omitempty"`
	Priority  *string    `json:"priority,omitempty"`
	Status    *string    `json:"status,omitempty"`
	Assignee  *string    `json:"assignee,omitempty"`
	ProjectID *string    `json:"project_id,omitempty"`
	TeamID    *string    `json:"team_id,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Label     *string    `json:"label,omitempty"`
	OwnerID   *string    `json:"owner_id,omitempty"`
}
