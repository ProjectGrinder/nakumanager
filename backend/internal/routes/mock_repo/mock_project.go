package mock

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockProjectRepo struct {
	mock.Mock
}

func (m *MockProjectRepo) CreateProject(ctx context.Context, data models.CreateProject) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockProjectRepo) DeleteProject(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepo) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	args := m.Called(ctx, id)
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

func (m *MockProjectRepo) IsProjectExists(ctx context.Context, projectID string) (bool, error) {
	args := m.Called(ctx, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepo) GetOwnerByProjectID(ctx context.Context, projectID string) (string, error) {
	args := m.Called(ctx, projectID)
	return args.String(0), args.Error(1)
}

func (m *MockProjectRepo) GetLeaderByProjectID(ctx context.Context, projectID string) (string, error) {
	args := m.Called(ctx, projectID)
	return args.String(0), args.Error(1)
}

func (m *MockProjectRepo) AddMemberToProject(ctx context.Context, projectID, userID string) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *MockProjectRepo) RemoveMemberFromProject(ctx context.Context, projectID, userID string) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}
