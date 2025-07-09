package repositories

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, data models.CreateProject) error
	DeleteProject(ctx context.Context, id string) error
	GetProjectByID(ctx context.Context, id string) (db.Project, error)
	GetProjectsByUserID(ctx context.Context, userID string) ([]db.Project, error)
	ListProjectMembers(ctx context.Context, projectID string) ([]db.User, error)
	ListProjectsByWorkspace(ctx context.Context, workspaceID string) ([]db.ListProjectsByWorkspaceRow, error)
	IsProjectExists(ctx context.Context, projectID string) (bool, error)
	GetOwnerByProjectID(ctx context.Context, projectID string) (string, error)
	GetLeaderByProjectID(ctx context.Context, projectID string) (string, error)
	AddMemberToProject(ctx context.Context, projectID, userID string) error
	RemoveMemberFromProject(ctx context.Context, projectID, userID string) error
}

type projectRepo struct {
	queries *db.Queries
}

func NewProjectRepository(q *db.Queries) ProjectRepository {
	return &projectRepo{queries: q}
}

func (r *projectRepo) CreateProject(ctx context.Context, data models.CreateProject) error {

	return r.queries.CreateProject(ctx, db.CreateProjectParams{
		ID:          data.ID,
		Name:        data.Name,
		LeaderID:    data.LeaderID,
		WorkspaceID: data.WorkspaceID,
		TeamID:      data.TeamID,
		Status:      data.Status,
		Priority:    data.Priority,
		StartDate:   data.StartDate,
		EndDate:     data.EndDate,
		Label:       data.Label,
		CreatedBy:   data.CreatedBy,
	})
}

func (r *projectRepo) DeleteProject(ctx context.Context, id string) error {
	return r.queries.DeleteProject(ctx, id)
}

func (r *projectRepo) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	return r.queries.GetProjectByID(ctx, id)
}

func (r *projectRepo) GetProjectsByUserID(ctx context.Context, userID string) ([]db.Project, error) {
	return r.queries.GetProjectsByUserID(ctx, db.GetProjectsByUserIDParams{UserID: userID, CreatedBy: userID})
}

func (r *projectRepo) ListProjectMembers(ctx context.Context, projectID string) ([]db.User, error) {
	return r.queries.ListProjectMembers(ctx, projectID)
}

func (r *projectRepo) ListProjectsByWorkspace(ctx context.Context, workspaceID string) ([]db.ListProjectsByWorkspaceRow, error) {
	return r.queries.ListProjectsByWorkspace(ctx, workspaceID)
}

func (r *projectRepo) IsProjectExists(ctx context.Context, projectID string) (bool, error) {
	exists, err := r.queries.IsProjectExists(ctx, projectID)
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *projectRepo) GetOwnerByProjectID(ctx context.Context, projectID string) (string, error) {
	return r.queries.GetOwnerByProjectID(ctx, projectID)
}

func (r *projectRepo) GetLeaderByProjectID(ctx context.Context, projectID string) (string, error) {
	result, _ := r.queries.GetLeaderByProjectID(ctx, projectID)

	leaderID, ok := result.(string)
	if !ok {
		return "", nil
	}
	return leaderID, nil

}

func (r *projectRepo) AddMemberToProject(ctx context.Context, projectID, userID string) error {
	return r.queries.AddMemberToProject(ctx, db.AddMemberToProjectParams{
		ProjectID: projectID,
		UserID:    userID,
	})
}

func (r *projectRepo) RemoveMemberFromProject(ctx context.Context, projectID, userID string) error {
	return r.queries.RemoveMemberFromProject(ctx, db.RemoveMemberFromProjectParams{
		ProjectID: projectID,
		UserID:    userID,
	})
}
