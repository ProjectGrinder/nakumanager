package model

import "log"

type Workspace struct {
	ID      string   `json:"id"`
	Name    string   `json:"name" validate:"required,min=3,max=50,alphanumunicode"`
	Members []string `json:"members"`
}

func (w *Workspace) Rename(newName string) {
	log.Println("reaname workspace")
}

func (w *Workspace) AddMember(user User) {
	log.Println("add member")
}

func (w *Workspace) RemoveMember(user User) {
	log.Println("remove member")
}
