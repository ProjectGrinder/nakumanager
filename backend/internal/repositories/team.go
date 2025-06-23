package repositories

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
)

type TeamRepository interface {
	AddMemberToTeam(ctx context.Context, data db.AddMemberToTeamParams) error
	RemoveMemberFromTeam(ctx context.Context, data db.RemoveMemberFromTeamParams) error
	CreateTeam(ctx context.Context, data db.CreateTeamParams) error
	DeleteTeam(ctx context.Context, id string) error
	GetTeamByID(ctx context.Context, id string) (db.Team, error)
	ListTeamMembers(ctx context.Context, teamID string) ([]db.User, error)
	ListTeams(ctx context.Context) ([]db.Team, error)
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

func (r *teamRepo) GetTeamByID(ctx context.Context, id string) (db.Team, error) {
	return r.queries.GetTeamByID(ctx, id)
}

func (r *teamRepo) ListTeamMembers(ctx context.Context, teamID string) ([]db.User, error) {
	return r.queries.ListTeamMembers(ctx, teamID)
}

func (r *teamRepo) ListTeams(ctx context.Context) ([]db.Team, error) {
	return r.queries.ListTeams(ctx)
}

