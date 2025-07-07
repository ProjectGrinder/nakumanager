package model

import "database/sql"

type CreateTeam struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspace_id"`
}

type AddMemberToTeam struct {
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

type RemoveMemberFromTeam struct {
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

type RenameTeam struct {
	TeamID string `json:"team_id"`
	Name   string `json:"name"`
}

type SetTeamLeader struct {
	TeamID   string         `json:"id"`
	LeaderID sql.NullString `json:"leader_id"`
}
