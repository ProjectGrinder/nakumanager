package repositories

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
)

type UserRepository interface {
	CreateUser(ctx context.Context, data db.CreateUserParams) error
	DeleteUser(ctx context.Context, id string) error
	GetUserByID(ctx context.Context, id string) (db.GetUserByIDRow, error)
	ListUsers(ctx context.Context) ([]db.ListUsersRow, error)
	UpdateEmail(ctx context.Context, data db.UpdateEmailParams) error
	UpdateRoles(ctx context.Context, data db.UpdateRolesParams) error
	UpdateUsername(ctx context.Context, data db.UpdateUsernameParams) error
	GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailWithoutPasswordRow, error)
	GetUserByEmailWithPassword(ctx context.Context, email string) (db.GetUserByEmailWithPasswordRow, error)
}

type userRepo struct {
	queries *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepo{queries: q}
}

func (r *userRepo) CreateUser(ctx context.Context, data db.CreateUserParams) error {
	err := r.queries.CreateUser(ctx, data)
	return err
}

func (r *userRepo) DeleteUser(ctx context.Context, id string) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (db.GetUserByIDRow, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *userRepo) ListUsers(ctx context.Context) ([]db.ListUsersRow, error) {
	return r.queries.ListUsers(ctx)
}

func (r *userRepo) UpdateEmail(ctx context.Context, data db.UpdateEmailParams) error {
	err := r.queries.UpdateEmail(ctx, data)
	return err
}

func (r *userRepo) UpdateRoles(ctx context.Context, data db.UpdateRolesParams) error {
	err := r.queries.UpdateRoles(ctx, data)
	return err
}

func (r *userRepo) UpdateUsername(ctx context.Context, data db.UpdateUsernameParams) error {
	err := r.queries.UpdateUsername(ctx, data)
	return err
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailWithoutPasswordRow, error) {
	return r.queries.GetUserByEmailWithoutPassword(ctx, email)
}

func (r *userRepo) GetUserByEmailWithPassword(ctx context.Context, email string) (db.GetUserByEmailWithPasswordRow, error) {
	return r.queries.GetUserByEmailWithPassword(ctx, email)
}
