package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

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

func (i *Issue) ChangeTitle(newTitle string) {
	i.Title = newTitle
}

func (i *Issue) ChangeContent(newContent string) {
	i.Content = newContent
}

func (i *Issue) ChangePriority(newPriority string) error {
	allowed := map[string]bool{"low": true, "medium": true, "high": true}
	if newPriority != "" && !allowed[newPriority] {
		return fmt.Errorf("invalid priority: %s", newPriority)
	}
	i.Priority = newPriority
	return nil
}

func (i *Issue) ChangeStatus(newStatus string) error {
	allowed := map[string]bool{"todo": true, "doing": true, "done": true}
	if newStatus != "" && !allowed[newStatus] {
		return fmt.Errorf("invalid status: %s", newStatus)
	}
	i.Status = newStatus
	return nil
}

func (i *Issue) AddAssignee(user User) {
	i.Assignee = append(i.Assignee, user.ID)
}

func (i *Issue) RemoveAssignee(user User) {
	newAssignees := i.Assignee[:0]
	for _, assignee := range i.Assignee {
		if assignee != user.ID {
			newAssignees = append(newAssignees, assignee)
		}
	}
	i.Assignee = newAssignees
}

func (i *Issue) AddProjectID(projectID uint) {
	i.ProjectID = projectID
}

func (i *Issue) AddTeamID(teamID uint) {
	i.TeamID = teamID
}

func (i *Issue) SetStartDate(t time.Time) {
	i.StartDate = &t
}

func (i *Issue) SetEndDate(t time.Time) {
	i.EndDate = &t
}

func (i *Issue) AddLabel(label string) {
	i.Label = label
}
