package model

import (
	"log"
)

type Team struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Members []User `json:"user"`
	Leader  User   `json:"leader"`
}

func (t *Team) Rename(newName string) {
	log.Println("rename team")
}

func (t *Team) AddMember(user User) {
	log.Println("add member")
}

func (t *Team) RemoveMember(user User) {
	log.Println("remove member")
}

func (t *Team) AssignLeader(user User) {
	log.Println("assign leader")
}
