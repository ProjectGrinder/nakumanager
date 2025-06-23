package model

import (
	"time"
)

type Issue struct {
	ID        string     `json:"id"`
	Title     string     `json:"title" validate:"required"`
	Content   string     `json:"content,omitempty"`
	Priority  string     `json:"priority" validate:"omitempty,oneof=low medium high"`
	Status    string     `json:"status" validate:"omitempty,oneof=todo doing done"`
	Assignee  []string   `json:"assignee,omitempty"`
	ProjectID uint       `json:"projectId,omitempty"`
	TeamID    uint       `json:"teamId" validate:"required"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Label     string     `json:"label,omitempty"`
	OwnerID   string     `json:"ownerId" validate:"required"`
}
