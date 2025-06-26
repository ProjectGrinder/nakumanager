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
	// Assert type for return values
	if proj, ok := args.Get(0).(db.Project); ok {
		return proj, args.Error(1)
	}
	return db.Project{}, args.Error(1)
}

func (m *MockProjectRepo) GetProjectsByUserID(ctx context.Context, userID string) ([]db.Project, error) {
	args := m.Called(ctx, userID)
	if projects, ok := args.Get(0).([]db.Project); ok {
		return projects, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectRepo) ListProjectMembers(ctx context.Context, projectID string) ([]db.User, error) {
	args := m.Called(ctx, projectID)
	if users, ok := args.Get(0).([]db.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectRepo) ListProjectsByWorkspace(ctx context.Context, workspaceID string) ([]db.ListProjectsByWorkspaceRow, error) {
	args := m.Called(ctx, workspaceID)
	if rows, ok := args.Get(0).([]db.ListProjectsByWorkspaceRow); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
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

func (m *MockProjectRepo) UpdateProjectName(ctx context.Context, projectID, name string) error {
	args := m.Called(ctx, projectID, name)
	return args.Error(0)
}

func (m *MockProjectRepo) UpdateProjectLeader(ctx context.Context, projectID, leaderID string) error {
	args := m.Called(ctx, projectID, leaderID)
	return args.Error(0)
}

func (m *MockProjectRepo) UpdateProjectWorkspace(ctx context.Context, projectID, workspaceID string) error {
	args := m.Called(ctx, projectID, workspaceID)
	return args.Error(0)
}
