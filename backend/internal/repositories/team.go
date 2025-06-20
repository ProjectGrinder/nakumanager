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