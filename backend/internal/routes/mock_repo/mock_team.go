package mock

import (
	"context"

	db "github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTeamRepo struct {
	mock.Mock
}

func (m *MockTeamRepo) AddMemberToTeam(ctx context.Context, data models.AddMemberToTeam) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTeamRepo) RemoveMemberFromTeam(ctx context.Context, data models.RemoveMemberFromTeam) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTeamRepo) CreateTeam(ctx context.Context, data models.CreateTeam) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTeamRepo) DeleteTeam(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTeamRepo) DeleteTeamFromTeamMembers(ctx context.Context, teamID string) error {
	args := m.Called(ctx, teamID)
	return args.Error(0)
}

func (m *MockTeamRepo) GetTeamByID(ctx context.Context, id string) (db.Team, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Team), args.Error(1)
}

func (m *MockTeamRepo) GetTeamsByUserID(ctx context.Context, userID string) ([]db.Team, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Team), args.Error(1)
}

func (m *MockTeamRepo) ListTeamMembers(ctx context.Context, teamID string) ([]db.ListTeamMembersRow, error) {
	args := m.Called(ctx, teamID)
	return args.Get(0).([]db.ListTeamMembersRow), args.Error(1)
}

func (m *MockTeamRepo) ListTeams(ctx context.Context) ([]db.Team, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.Team), args.Error(1)
}

func (m *MockTeamRepo) GetOwnerByTeamID(ctx context.Context, teamID string) (string, error) {
	args := m.Called(ctx, teamID)
	return args.String(0), args.Error(1)
}

func (m *MockTeamRepo) GetLeaderByTeamID(ctx context.Context, teamID string) (string, error) {
	args := m.Called(ctx, teamID)
	return args.String(0), args.Error(1)
}

func (m *MockTeamRepo) IsMemberInTeam(ctx context.Context, teamID, userID string) (bool, error) {
	args := m.Called(ctx, teamID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTeamRepo) IsTeamExists(ctx context.Context, teamID string) (bool, error) {
	args := m.Called(ctx, teamID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTeamRepo) RenameTeam(ctx context.Context, data models.RenameTeam) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTeamRepo) SetLeaderToTeam(ctx context.Context, data models.SetTeamLeader) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}
