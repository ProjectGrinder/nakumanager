package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nack098/nakumanager/internal/db"
)

type TeamRepository interface {
	AddMemberToTeam(ctx context.Context, data db.AddMemberToTeamParams) error
	RemoveMemberFromTeam(ctx context.Context, data db.RemoveMemberFromTeamParams) error
	CreateTeam(ctx context.Context, data db.CreateTeamParams) error
	DeleteTeam(ctx context.Context, id string) error
	DeleteTeamFromTeamMembers(ctx context.Context, teamID string) error
	GetTeamByID(ctx context.Context, id string) (db.Team, error)
	GetTeamsByUserID(ctx context.Context, userID string) ([]db.Team, error)
	ListTeamMembers(ctx context.Context, teamID string) ([]db.ListTeamMembersRow, error)
	ListTeams(ctx context.Context) ([]db.Team, error)
	GetOwnerByTeamID(ctx context.Context, teamID string) (string, error)
	GetLeaderByTeamID(ctx context.Context, userID string) (string, error)
	IsMemberInTeam(ctx context.Context, teamID, userID string) (bool, error)
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

func (r *teamRepo) CreateTeam(ctx context.Context, data db.CreateTeamParams) error {
	return r.queries.CreateTeam(ctx, data)
}

func (r *teamRepo) DeleteTeam(ctx context.Context, id string) error {
	return r.queries.DeleteTeam(ctx, id)
}

func (r *teamRepo) DeleteTeamFromTeamMembers(ctx context.Context, teamID string) error {
	return r.queries.DeleteTeamFromTeamMembers(ctx, teamID)
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
