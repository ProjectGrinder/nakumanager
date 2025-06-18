package model

import (
	"time"
	"log"
)

type Issue struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Piority   string `json:"priority"`
	Status    string `json:"status"`
	Assignee  []User `json:"assignee"`
	ProjectID uint
	TeamID    uint
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Label     string    `json:"label"`
}

func (i *Issue) changeTitle(newTitle string) {
	log.Println("change title")
}

func (i *Issue) changeContent(newContent string) {
	log.Println("change content")
}

func (i *Issue) changePriority(newPriority string) {
	log.Println("change priority")
}

func (i *Issue) changeStatus(newStatus string) {
	log.Println("change status")
}

func (i *Issue) addAssignee(user User) {
	log.Println("add assignee")
}

func (i *Issue) removeAssignee(user User) {
	log.Println("remove assignee")
}

func (i *Issue) addProjectID(projectID uint) {
	log.Println("add project id")
}

func (i *Issue) addTeamID(teamID uint) {
	log.Println("add team id")
}
