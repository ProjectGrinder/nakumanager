package model

import (
	"time"
)

type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Piority     string `json:"priority"`
	WorkspaceID uint
	Leader      User      `json:"leader"`
	Members     []User    `json:"members"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Label       string    `json:"label"`
}
