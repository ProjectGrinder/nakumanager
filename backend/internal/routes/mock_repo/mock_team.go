package mock

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) AddMemberToTeam(ctx context.Context, data db.AddMemberToTeamParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTeamRepository) RemoveMemberFromTeam(ctx context.Context, data db.RemoveMemberFromTeamParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTeamRepository) CreateTeam(ctx context.Context, data models.CreateTeam) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTeamRepository) DeleteTeam(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTeamRepository) GetTeamByID(ctx context.Context, id string) (db.Team, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Team), args.Error(1)
}

func (m *MockTeamRepository) GetTeamsByUserID(ctx context.Context, userID string) ([]db.Team, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Team), args.Error(1)
}

func (m *MockTeamRepository) GetOwnerByTeamID(ctx context.Context, teamID string) (string, error) {
	args := m.Called(ctx, teamID)
	return args.String(0), args.Error(1)
}

func (m *MockTeamRepository) GetLeaderByTeamID(ctx context.Context, teamID string) (string, error) {
	args := m.Called(ctx, teamID)
	return args.String(0), args.Error(1)
}

func (m *MockTeamRepository) IsMemberInTeam(ctx context.Context, teamID, userID string) (bool, error) {
	args := m.Called(ctx, teamID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTeamRepository) IsTeamExists(ctx context.Context, teamID string) (bool, error) {
	args := m.Called(ctx, teamID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTeamRepository) RenameTeam(ctx context.Context, data db.RenameTeamParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTeamRepository) SetLeaderToTeam(ctx context.Context, data db.SetLeaderToTeamParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}
