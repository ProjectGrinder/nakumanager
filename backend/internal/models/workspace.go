package model

type CreateWorkspace struct {
	ID      string   `json:"id"`
	Name    string   `json:"name" validate:"required"`
	Members []string `json:"members"`
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
