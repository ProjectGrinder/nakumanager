package mock

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
	"github.com/stretchr/testify/mock"
)

type MockIssueRepo struct {
	mock.Mock
}

func (m *MockIssueRepo) AddAssigneeToIssue(ctx context.Context, data db.AddAssigneeToIssueParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockIssueRepo) CreateIssue(ctx context.Context, data db.CreateIssueParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockIssueRepo) DeleteIssue(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockIssueRepo) GetIssueByID(ctx context.Context, id string) (db.Issue, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Issue), args.Error(1)
}

func (m *MockIssueRepo) ListAssigneesByIssueID(ctx context.Context, issueID string) ([]db.User, error) {
	args := m.Called(ctx, issueID)
	return args.Get(0).([]db.User), args.Error(1)
}

func (m *MockIssueRepo) ListIssuesByTeamID(ctx context.Context, teamID string) ([]db.Issue, error) {
	args := m.Called(ctx, teamID)
	return args.Get(0).([]db.Issue), args.Error(1)
}

func (m *MockIssueRepo) RemoveAssigneeFromIssue(ctx context.Context, data db.RemoveAssigneeFromIssueParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockIssueRepo) GetIssueByUserID(ctx context.Context, userID string) ([]db.Issue, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Issue), args.Error(1)
}

func (m *MockIssueRepo) UpdateIssue(ctx context.Context, params db.UpdateIssueParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}
