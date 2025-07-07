package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
)

type TeamRepository interface {
	AddMemberToTeam(ctx context.Context, data models.AddMemberToTeam) error
	RemoveMemberFromTeam(ctx context.Context, data models.RemoveMemberFromTeam) error
	CreateTeam(ctx context.Context, data models.CreateTeam) error
	DeleteTeam(ctx context.Context, id string) error
	GetTeamByID(ctx context.Context, id string) (db.Team, error)
	GetTeamsByUserID(ctx context.Context, userID string) ([]db.Team, error)
	ListTeamMembers(ctx context.Context, teamID string) ([]db.ListTeamMembersRow, error)
	ListTeams(ctx context.Context) ([]db.Team, error)
	GetOwnerByTeamID(ctx context.Context, teamID string) (string, error)
	GetLeaderByTeamID(ctx context.Context, userID string) (string, error)
	IsMemberInTeam(ctx context.Context, teamID, userID string) (bool, error)
	IsTeamExists(ctx context.Context, teamID string) (bool, error)
	RenameTeam(ctx context.Context, data models.RenameTeam) error
	SetLeaderToTeam(ctx context.Context, data models.SetTeamLeader) error
}

type teamRepo struct {
	queries *db.Queries
}

func NewTeamRepository(q *db.Queries) TeamRepository {
	return &teamRepo{queries: q}
}

func (r *teamRepo) AddMemberToTeam(ctx context.Context, data models.AddMemberToTeam) error {
	model := db.AddMemberToTeamParams{
		TeamID: data.TeamID,
		UserID: data.UserID,
	}
	return r.queries.AddMemberToTeam(ctx, model)
}

func (r *teamRepo) RemoveMemberFromTeam(ctx context.Context, data models.RemoveMemberFromTeam) error {
	model := db.RemoveMemberFromTeamParams{
		TeamID: data.TeamID,
		UserID: data.UserID,
	}
	return r.queries.RemoveMemberFromTeam(ctx, model)
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

func (r *teamRepo) ListTeamMembers(ctx context.Context, teamID string) ([]db.ListTeamMembersRow, error) {
	return r.queries.ListTeamMembers(ctx, teamID)
}

func (r *teamRepo) ListTeams(ctx context.Context) ([]db.Team, error) {
	return r.queries.ListTeams(ctx)
}

func (r *teamRepo) GetOwnerByTeamID(ctx context.Context, teamID string) (string, error) {
	owner, err := r.queries.GetOwnerByTeamID(ctx, teamID)
	return owner, err
}

func (r *teamRepo) GetLeaderByTeamID(ctx context.Context, teamID string) (string, error) {
	leader, err := r.queries.GetLeaderByTeamID(ctx, teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("team %s not found", teamID)
		}
		return "", err
	}
	if !leader.Valid {
		return "", nil
	}
	return leader.String, nil
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

func (r *teamRepo) RenameTeam(ctx context.Context, data models.RenameTeam) error {
	model := db.RenameTeamParams{ID: data.TeamID, Name: data.Name}
	return r.queries.RenameTeam(ctx, model)
}

func (r *teamRepo) SetLeaderToTeam(ctx context.Context, data models.SetTeamLeader) error {
	return r.queries.SetLeaderToTeam(ctx, db.SetLeaderToTeamParams{
		LeaderID: data.LeaderID,
		ID:       data.TeamID,
	})
}
