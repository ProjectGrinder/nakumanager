package mock

import (
	"context"

	db "github.com/nack098/nakumanager/internal/db"
	"github.com/stretchr/testify/mock"
)

type MockWorkspaceRepo struct {
	mock.Mock
}

func (m *MockWorkspaceRepo) CreateWorkspace(ctx context.Context, id string, name string, ownerID string) error {
	args := m.Called(ctx, id, name, ownerID)
	return args.Error(0)
}

func (m *MockWorkspaceRepo) GetWorkspaceByID(ctx context.Context, id string) (db.Workspace, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Workspace), args.Error(1)
}

func (m *MockWorkspaceRepo) DeleteWorkspace(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWorkspaceRepo) AddMemberToWorkspace(ctx context.Context, workspaceID, userID string) error {
	args := m.Called(ctx, workspaceID, userID)
	return args.Error(0)
}

func (m *MockWorkspaceRepo) RemoveMemberFromWorkspace(ctx context.Context, workspaceID, userID string) error {
	args := m.Called(ctx, workspaceID, userID)
	return args.Error(0)
}

func (m *MockWorkspaceRepo) RenameWorkspace(ctx context.Context, id string, newName string) error {
	args := m.Called(ctx, id, newName)
	return args.Error(0)
}

func (m *MockWorkspaceRepo) ListWorkspacesWithMembersByUserID(ctx context.Context, userID string) ([]db.ListWorkspacesWithMembersByUserIDRow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.ListWorkspacesWithMembersByUserIDRow), args.Error(1)
}
