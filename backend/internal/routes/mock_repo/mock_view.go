package mock

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
	"github.com/stretchr/testify/mock"
)

type MockViewRepo struct {
	mock.Mock
}

func (m *MockViewRepo) AddGroupByToView(ctx context.Context, data db.AddGroupByToViewParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockViewRepo) AddIssueToView(ctx context.Context, data db.AddIssueToViewParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockViewRepo) CreateView(ctx context.Context, data db.CreateViewParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockViewRepo) DeleteView(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockViewRepo) GetViewByID(ctx context.Context, id string) ([]db.View, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]db.View), args.Error(1)
}

func (m *MockViewRepo) ListGroupByViewID(ctx context.Context, viewID string) ([]string, error) {
	args := m.Called(ctx, viewID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockViewRepo) ListIssuesByViewID(ctx context.Context, viewID string) ([]db.Issue, error) {
	args := m.Called(ctx, viewID)
	return args.Get(0).([]db.Issue), args.Error(1)
}

func (m *MockViewRepo) ListViewsByUser(ctx context.Context, userID string) ([]db.View, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.View), args.Error(1)
}

func (m *MockViewRepo) UpdateViewName(ctx context.Context, id string, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockViewRepo) RemoveGroupByFromView(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockViewRepo) RemoveIssueFromView(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockViewRepo) ListIssuesByGroupFilters(ctx context.Context, teamID string, filters map[string]string) ([]db.Issue, error) {
	args := m.Called(ctx, teamID, filters)
	return args.Get(0).([]db.Issue), args.Error(1)
}

func (m *MockViewRepo) GetGroupedIssues(ctx context.Context, teamID string, groupBy []string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, teamID, groupBy)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}
