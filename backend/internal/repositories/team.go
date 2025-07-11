package repositories

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
)

type TeamRepository interface {
	AddMemberToTeam(ctx context.Context, data db.AddMemberToTeamParams) error
	RemoveMemberFromTeam(ctx context.Context, data db.RemoveMemberFromTeamParams) error
	CreateTeam(ctx context.Context, data models.CreateTeam) error
	DeleteTeam(ctx context.Context, id string) error
	GetTeamByID(ctx context.Context, id string) (db.Team, error)
	GetTeamsByUserID(ctx context.Context, userID string) ([]db.Team, error)
	GetOwnerByTeamID(ctx context.Context, teamID string) (string, error)
	GetLeaderByTeamID(ctx context.Context, userID string) (string, error)
	IsMemberInTeam(ctx context.Context, teamID, userID string) (bool, error)
	IsTeamExists(ctx context.Context, teamID string) (bool, error)
	RenameTeam(ctx context.Context, data db.RenameTeamParams) error
	SetLeaderToTeam(ctx context.Context, data db.SetLeaderToTeamParams) error
}

type teamRepo struct {
	queries *db.Queries
}

func NewTeamRepository(q *db.Queries) TeamRepository {
	return &teamRepo{queries: q}
}

func (r *teamRepo) AddMemberToTeam(ctx context.Context, data db.AddMemberToTeamParams) error {
	return r.queries.AddMemberToTeam(ctx, data)
}

func (r *teamRepo) RemoveMemberFromTeam(ctx context.Context, data db.RemoveMemberFromTeamParams) error {
	return r.queries.RemoveMemberFromTeam(ctx, data)
}

func (r *teamRepo) CreateTeam(ctx context.Context, data models.CreateTeam) error {
	model := db.CreateTeamParams{
		ID:          data.ID,
		Name:        data.Name,
		WorkspaceID: data.WorkspaceID,
	}
	return r.queries.CreateTeam(ctx, model)
}

func (r *teamRepo) DeleteTeam(ctx context.Context, id string) error {
	return r.queries.DeleteTeam(ctx, id)
}

func (r *teamRepo) GetTeamByID(ctx context.Context, id string) (db.Team, error) {
	return r.queries.GetTeamByID(ctx, id)
}

func (r *teamRepo) GetTeamsByUserID(ctx context.Context, userID string) ([]db.Team, error) {
	return r.queries.GetTeamsByUserID(ctx, userID)
}

func (r *teamRepo) GetOwnerByTeamID(ctx context.Context, teamID string) (string, error) {
	owner, err := r.queries.GetOwnerByTeamID(ctx, teamID)
	return owner, err
}

func (r *teamRepo) GetLeaderByTeamID(ctx context.Context, teamID string) (string, error) {
	// leader, err := r.queries.GetLeaderByTeamID(ctx, teamID)
	// return leader, err
	return "", nil
}

func (r *teamRepo) IsMemberInTeam(ctx context.Context, teamID, userID string) (bool, error) {
	count, err := r.queries.IsMemberInTeam(ctx, db.IsMemberInTeamParams{
		TeamID: teamID,
		UserID: userID,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *teamRepo) IsTeamExists(ctx context.Context, teamID string) (bool, error) {
	count, err := r.queries.IsTeamExists(ctx, teamID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *teamRepo) RenameTeam(ctx context.Context, data db.RenameTeamParams) error {
	return r.queries.RenameTeam(ctx, data)
}

func (r *teamRepo) SetLeaderToTeam(ctx context.Context, data db.SetLeaderToTeamParams) error {
	return r.queries.SetLeaderToTeam(ctx, data)
}
