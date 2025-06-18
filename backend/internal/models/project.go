package model

import (
	"log"
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

func (p *Project) changeStatus(newStatus string) {
	log.Println("change status")
}

func (p *Project) changePriority(newPriority string) {
	log.Println("change priority")
}

func (p *Project) changeStartDate(newStartDate time.Time) {
	log.Println("change start date")
}

func (p *Project) changeEndDate(newEndDate time.Time) {
	log.Println("change end date")
}

func (p *Project) changeLabel(newLabel string) {
	log.Println("change label")
}

func (p *Project) changeLeader(newLeader User) {
	log.Println("change leader")
}

func (p *Project) addMember(user User) {
	log.Println("add member")
}

func (p *Project) removeMember(user User) {
	log.Println("remove member")
}

func (p *Project) addWorkSpaceID(newWorkSpaceID uint) {
	log.Println("new workspace id")
}
