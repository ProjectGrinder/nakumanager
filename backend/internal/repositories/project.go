package repositories

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
)

type ProjectRepository interface {
	AddMemberToProject(ctx context.Context, data db.AddMemberToProjectParams) error
	CreateProject(ctx context.Context, data db.CreateProjectParams) error
	DeleteProject(ctx context.Context, id string) error
	GetProjectByID(ctx context.Context, id string) (db.Project, error)
	GetProjectsByUserID(ctx context.Context, userID string) ([]db.Project, error)
	ListProjectMembers(ctx context.Context, projectID string) ([]db.User, error)
	ListProjectsByWorkspace(ctx context.Context, workspaceID string) ([]db.ListProjectsByWorkspaceRow, error)
	RemoveMemberFromProject(ctx context.Context, data db.RemoveMemberFromProjectParams) error
	UpdateProject(ctx context.Context, data db.UpdateProjectParams) error
	IsProjectExists(ctx context.Context, projectID string) (bool, error)
}

type projectRepo struct {
	queries *db.Queries
}

func NewProjectRepository(q *db.Queries) ProjectRepository {
	return &projectRepo{queries: q}
}

func (r *projectRepo) AddMemberToProject(ctx context.Context, data db.AddMemberToProjectParams) error {
	return r.queries.AddMemberToProject(ctx, data)
}

func (r *projectRepo) CreateProject(ctx context.Context, data db.CreateProjectParams) error {
	return r.queries.CreateProject(ctx, data)
}

func (r *projectRepo) DeleteProject(ctx context.Context, id string) error {
	return r.queries.DeleteProject(ctx, id)
}

func (r *projectRepo) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	return r.queries.GetProjectByID(ctx, id)
}

func (r *projectRepo) GetProjectsByUserID(ctx context.Context, userID string) ([]db.Project, error) {
	return r.queries.GetProjectsByUserID(ctx, userID)
}

func (r *projectRepo) ListProjectMembers(ctx context.Context, projectID string) ([]db.User, error) {
	return r.queries.ListProjectMembers(ctx, projectID)
}

func (r *projectRepo) ListProjectsByWorkspace(ctx context.Context, workspaceID string) ([]db.ListProjectsByWorkspaceRow, error) {
	return r.queries.ListProjectsByWorkspace(ctx, workspaceID)
}

func (r *projectRepo) RemoveMemberFromProject(ctx context.Context, data db.RemoveMemberFromProjectParams) error {
	return r.queries.RemoveMemberFromProject(ctx, data)
}

func (r *projectRepo) UpdateProject(ctx context.Context, data db.UpdateProjectParams) error {
	return r.queries.UpdateProject(ctx, data)
}

func (r *projectRepo) IsProjectExists(ctx context.Context, projectID string) (bool, error) {
	exists, err := r.queries.IsProjectExists(ctx, projectID)
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
