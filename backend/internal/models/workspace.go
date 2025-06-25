package model

type Workspace struct {
	ID      string   `json:"id"`
	Name    string   `json:"name" validate:"required"`
	Members []string `json:"members"`
}

type EditWorkspace struct {
	WorkspaceID  string `json:"workspace_id"`
	Name         string `json:"name"`
	AddMember    string `json:"add_member"`
	RemoveMember string `json:"remove_member"`
}
