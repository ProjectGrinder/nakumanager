package mock

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
	"github.com/stretchr/testify/mock"
)

type MockProjectRepo struct {
	mock.Mock
}

func (m *MockProjectRepo) AddMemberToProject(ctx context.Context, data db.AddMemberToProjectParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockProjectRepo) CreateProject(ctx context.Context, data db.CreateProjectParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockProjectRepo) DeleteProject(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepo) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Project), args.Error(1)
}

func (m *MockProjectRepo) GetProjectsByUserID(ctx context.Context, userID string) ([]db.Project, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Project), args.Error(1)
}

func (m *MockProjectRepo) ListProjectMembers(ctx context.Context, projectID string) ([]db.User, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]db.User), args.Error(1)
}

func (m *MockProjectRepo) ListProjectsByWorkspace(ctx context.Context, workspaceID string) ([]db.ListProjectsByWorkspaceRow, error) {
	args := m.Called(ctx, workspaceID)
	return args.Get(0).([]db.ListProjectsByWorkspaceRow), args.Error(1)
}

func (m *MockProjectRepo) RemoveMemberFromProject(ctx context.Context, data db.RemoveMemberFromProjectParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockProjectRepo) UpdateProject(ctx context.Context, data db.UpdateProjectParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockProjectRepo) IsProjectExists(ctx context.Context, projectID string) (bool, error) {
	args := m.Called(ctx, projectID)
	return args.Bool(0), args.Error(1)
}
