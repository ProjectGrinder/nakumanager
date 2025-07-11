package model

type CreateTeam struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"required"`
	WorkspaceID string `json:"workspace_id" validate:"required"`
}

type UpdateTeamRequest struct {
	Name          *string   `json:"Name"`
	AddMembers    *[]string `json:"add_members"`
	RemoveMembers *[]string `json:"remove_members"`
	NewLeaderID   *string   `json:"new_leader_id"`
}
