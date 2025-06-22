package repositories

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
)

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, id string, name string, ownerID string) error
	GetWorkspaceByID(ctx context.Context, id string) (db.Workspace, error)
	GetWorkspaceByUserID(ctx context.Context, userID string) ([]db.Workspace, error)
	DeleteWorkspace(ctx context.Context, id string) error
	ListWorkspaceMembers(ctx context.Context, workspaceID string) ([]db.User, error)
	AddMemberToWorkspace(ctx context.Context, workspaceID, userID string) error
	RemoveMemberFromWorkspace(ctx context.Context, workspaceID, userID string) error
	RenameWorkspace(ctx context.Context, id string, newName string) error
	ListWorkspacesWithMembersByUserID(ctx context.Context, userID string) ([]db.ListWorkspacesWithMembersByUserIDRow, error)
}

type workspaceRepo struct {
	queries *db.Queries
}

func NewWorkspaceRepository(q *db.Queries) WorkspaceRepository {
	return &workspaceRepo{queries: q}
}

func (r *workspaceRepo) CreateWorkspace(ctx context.Context, id string, name string, ownerID string) error {
	return r.queries.CreateWorkspace(ctx, db.CreateWorkspaceParams{
		ID:      id,
		Name:    name,
		OwnerID: ownerID,
	})
}

func (r *workspaceRepo) GetWorkspaceByID(ctx context.Context, id string) (db.Workspace, error) {
	return r.queries.GetWorkspaceByID(ctx, id)
}

func (r *workspaceRepo) GetWorkspaceByUserID(ctx context.Context, userID string) ([]db.Workspace, error) {
	workspaces, err := r.queries.GetWorkspaceByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (r *workspaceRepo) DeleteWorkspace(ctx context.Context, id string) error {
	return r.queries.DeleteWorkspace(ctx, id)
}

func (r *workspaceRepo) ListWorkspaceMembers(ctx context.Context, workspaceID string) ([]db.User, error) {
	return r.queries.ListWorkspaceMembers(ctx, workspaceID)
}

func (r *workspaceRepo) AddMemberToWorkspace(ctx context.Context, workspaceID, userID string) error {
	return r.queries.AddMemberToWorkspace(ctx, db.AddMemberToWorkspaceParams{
		WorkspaceID: workspaceID,
		UserID:      userID,
	})
}

func (r *workspaceRepo) RemoveMemberFromWorkspace(ctx context.Context, workspaceID, userID string) error {
	return r.queries.RemoveMemberFromWorkspace(ctx, db.RemoveMemberFromWorkspaceParams{
		WorkspaceID: workspaceID,
		UserID:      userID,
	})
}

func (r *workspaceRepo) RenameWorkspace(ctx context.Context, id string, newName string) error {
	return r.queries.RenameWorkspace(ctx, db.RenameWorkspaceParams{
		Name: newName,
		ID:   id,
	})
}

func (r *workspaceRepo) ListWorkspacesWithMembersByUserID(ctx context.Context, userID string) ([]db.ListWorkspacesWithMembersByUserIDRow, error) {
	params := db.ListWorkspacesWithMembersByUserIDParams{
		OwnerID: userID,
		UserID:  userID,
	}
	return r.queries.ListWorkspacesWithMembersByUserID(ctx, params)
}
