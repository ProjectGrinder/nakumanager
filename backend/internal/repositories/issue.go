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
	ListIssueAssignees(ctx context.Context, issueID string) ([]db.User, error)
	ListIssuesByProjectID(ctx context.Context, projectID string) ([]db.Issue, error)
	ListIssuesByTeamID(ctx context.Context, teamID string) ([]db.Issue, error)
	RemoveAssigneeFromIssue(ctx context.Context, data db.RemoveAssigneeFromIssueParams) error
	UpdateIssue(ctx context.Context, data db.UpdateIssueParams) error
	UpdateIssueAssignees(ctx context.Context, data db.UpdateIssueAssigneesParams) error
	UpdateIssueStatus(ctx context.Context, data db.UpdateIssueStatusParams) error
}