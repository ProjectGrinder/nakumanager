package model

import "log"

type Workspace struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Members []User `json:"user"`
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
