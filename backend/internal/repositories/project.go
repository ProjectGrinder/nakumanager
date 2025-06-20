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
	ListProjectMembers(ctx context.Context, projectID string) ([]db.User, error)
	ListProjectsByWorkspace(ctx context.Context, workspaceID string) ([]db.Project, error)
	RemoveMemberFromProject(ctx context.Context, data db.RemoveMemberFromProjectParams) error
	UpdateProject(ctx context.Context, data db.UpdateProjectParams) error
}
