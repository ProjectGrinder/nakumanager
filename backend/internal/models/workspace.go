package model

type CreateWorkspace struct {
	ID      string   `json:"id"`
	Name    string   `json:"name" validate:"required"`
	Members []string `json:"members"`
}

type UpdateWorkspaceRequest struct {
	Name          *string   `json:"name"`
	AddMembers    *[]string `json:"add_members"`
	RemoveMembers *[]string `json:"remove_members"`
}
