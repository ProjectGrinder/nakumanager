package model

type CreateWorkspace struct {
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

type AddMemberRequest struct {
	MemberID string `json:"member_id" validate:"required"`
}

type RemoveMemberRequest struct {
	MemberID string `json:"member_id" validate:"required"`
}

type RenameWorkSpaceRequest struct {
	Name string `json:"name" validate:"required"`
}
