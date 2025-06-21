package repositories

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
)

type IssueRepository interface {
	AddAssigneeToIssue(ctx context.Context, data db.AddAssigneeToIssueParams) error
	CreateIssue(ctx context.Context, data db.CreateIssueParams) error
	DeleteIssue(ctx context.Context, id string) error
	GetIssueByID(ctx context.Context, id string) (db.Issue, error)
	ListAssigneesByIssueID(ctx context.Context, issueID string) ([]db.User, error)
	ListIssuesByProjectID(ctx context.Context, projectID string) ([]db.Issue, error)
	ListIssuesByTeamID(ctx context.Context, teamID string) ([]db.Issue, error)
	RemoveAssigneeFromIssue(ctx context.Context, data db.RemoveAssigneeFromIssueParams) error
	UpdateIssue(ctx context.Context, data db.UpdateIssueParams) error
	UpdateIssueAssignees(ctx context.Context, data db.UpdateIssueAssigneesParams) error
	UpdateIssueStatus(ctx context.Context, data db.UpdateIssueStatusParams) error
}

type issueRepo struct {
	queries *db.Queries
}

func NewIssueRepository(q *db.Queries) IssueRepository {
	return &issueRepo{queries: q}
}

func (r *issueRepo) AddAssigneeToIssue(ctx context.Context, data db.AddAssigneeToIssueParams) error {
	return r.queries.AddAssigneeToIssue(ctx, data)
}

func (r *issueRepo) CreateIssue(ctx context.Context, data db.CreateIssueParams) error {
	return r.queries.CreateIssue(ctx, data)
}

func (r *issueRepo) DeleteIssue(ctx context.Context, id string) error {
	return r.queries.DeleteIssue(ctx, id)
}

func (r *issueRepo) GetIssueByID(ctx context.Context, id string) (db.Issue, error) {
	return r.queries.GetIssueByID(ctx, id)
}

func (r *issueRepo) ListAssigneesByIssueID(ctx context.Context, issueID string) ([]db.User, error) {
	return r.queries.ListAssigneesByIssueID(ctx, issueID)
}

func (r *issueRepo) ListIssuesByProjectID(ctx context.Context, projectID string) ([]db.Issue, error) {
	return r.queries.ListIssuesByProjectID(ctx, projectID)
}

func (r *issueRepo) ListIssuesByTeamID(ctx context.Context, teamID string) ([]db.Issue, error) {
	return r.queries.ListIssuesByTeamID(ctx, teamID)
}

func (r *issueRepo) RemoveAssigneeFromIssue(ctx context.Context, data db.RemoveAssigneeFromIssueParams) error {
	return r.queries.RemoveAssigneeFromIssue(ctx, data)
}

func (r *issueRepo) UpdateIssue(ctx context.Context, data db.UpdateIssueParams) error {
	return r.queries.UpdateIssue(ctx, data)
}

func (r *issueRepo) UpdateIssueAssignees(ctx context.Context, data db.UpdateIssueAssigneesParams) error {
	return r.queries.UpdateIssueAssignees(ctx, data)
}

func (r *issueRepo) UpdateIssueStatus(ctx context.Context, data db.UpdateIssueStatusParams) error {
	return r.queries.UpdateIssueStatus(ctx, data)
}
