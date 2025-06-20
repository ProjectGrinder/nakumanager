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