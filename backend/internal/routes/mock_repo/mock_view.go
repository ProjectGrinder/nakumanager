package mock

import (
	"context"
	"database/sql"

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

func (m *MockViewRepo) AddIssueToViewTx(ctx context.Context, tx *sql.Tx, params db.AddIssueToViewParams) error {
	args := m.Called(ctx, tx, params)
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

	if data := args.Get(0); data != nil {
		return data.([]db.View), args.Error(1)
	}
	return nil, args.Error(1)
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

func (m *MockViewRepo) ListViewByTeamID(ctx context.Context, teamID string) ([]db.View, error) {
	args := m.Called(ctx, teamID)

	if data := args.Get(0); data != nil {
		return data.([]db.View), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockViewRepo) GetViewsByGroupBys(ctx context.Context, groupBys []string) ([]db.View, error) {
	args := m.Called(ctx, groupBys)

	if data := args.Get(0); data != nil {
		return data.([]db.View), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockViewRepo) UpdateViewTeamID(ctx context.Context, id string, teamID string) error {
	args := m.Called(ctx, id, teamID)
	return args.Error(0)
}

func (m *MockViewRepo) GetTeamIDByViewID(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}
