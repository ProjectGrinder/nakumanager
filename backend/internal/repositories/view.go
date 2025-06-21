package repositories

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
)

type ViewRepository interface {
	AddGroupByToView(ctx context.Context, data db.AddGroupByToViewParams) error
	AddIssueToView(ctx context.Context, data db.AddIssueToViewParams) error
	CreateView(ctx context.Context, data db.CreateViewParams) error
	DeleteView(ctx context.Context, id string) error
	GetViewByID(ctx context.Context, id string) (db.View, error)
	ListGroupBysByViewID(ctx context.Context, viewID string) ([]string, error)
	ListIssuesByViewID(ctx context.Context, viewID string) ([]db.Issue, error)
	ListViewsByUser(ctx context.Context, userID string) ([]db.View, error)
	RemoveGroupByFromView(ctx context.Context, data db.RemoveGroupByFromViewParams) error
	RemoveIssueFromView(ctx context.Context, data db.RemoveIssueFromViewParams) error
}

type viewRepo struct {
	db *db.Queries
}

func NewViewRepository(q *db.Queries) ViewRepository {
	return &viewRepo{db: q}
}

func (r *viewRepo) AddGroupByToView(ctx context.Context, data db.AddGroupByToViewParams) error {
	return r.db.AddGroupByToView(ctx, data)
}

func (r *viewRepo) AddIssueToView(ctx context.Context, data db.AddIssueToViewParams) error {
	return r.db.AddIssueToView(ctx, data)
}

func (r *viewRepo) CreateView(ctx context.Context, data db.CreateViewParams) error {
	return r.db.CreateView(ctx, data)
}

func (r *viewRepo) DeleteView(ctx context.Context, id string) error {
	return r.db.DeleteView(ctx, id)
}

func (r *viewRepo) GetViewByID(ctx context.Context, id string) (db.View, error) {
	return r.db.GetViewByID(ctx, id)
}

func (r *viewRepo) ListGroupBysByViewID(ctx context.Context, viewID string) ([]string, error) {
	return r.db.ListGroupBysByViewID(ctx, viewID)
}

func (r *viewRepo) ListIssuesByViewID(ctx context.Context, viewID string) ([]db.Issue, error) {
	return r.db.ListIssuesByViewID(ctx, viewID)
}

func (r *viewRepo) ListViewsByUser(ctx context.Context, userID string) ([]db.View, error) {
	return r.db.ListViewsByUser(ctx, userID)
}

func (r *viewRepo) RemoveGroupByFromView(ctx context.Context, data db.RemoveGroupByFromViewParams) error {
	return r.db.RemoveGroupByFromView(ctx, data)
}

func (r *viewRepo) RemoveIssueFromView(ctx context.Context, data db.RemoveIssueFromViewParams) error {
	return r.db.RemoveIssueFromView(ctx, data)
}
