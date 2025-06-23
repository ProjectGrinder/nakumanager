package mock

import (
	"context"

	"github.com/nack098/nakumanager/internal/db"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, data db.CreateUserParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id string) (db.GetUserByIDRow, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(db.GetUserByIDRow); ok {
		return user, args.Error(1)
	}
	return db.GetUserByIDRow{}, args.Error(1)
}

func (m *MockUserRepo) ListUsers(ctx context.Context) ([]db.ListUsersRow, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]db.ListUsersRow); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) UpdateEmail(ctx context.Context, data db.UpdateEmailParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateRoles(ctx context.Context, data db.UpdateRolesParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateUsername(ctx context.Context, data db.UpdateUsernameParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailWithoutPasswordRow, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(db.GetUserByEmailWithoutPasswordRow); ok {
		return user, args.Error(1)
	}
	return db.GetUserByEmailWithoutPasswordRow{}, args.Error(1)
}

func (m *MockUserRepo) GetUserByEmailWithPassword(ctx context.Context, email string) (db.GetUserByEmailWithPasswordRow, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(db.GetUserByEmailWithPasswordRow); ok {
		return user, args.Error(1)
	}
	return db.GetUserByEmailWithPasswordRow{}, args.Error(1)
}
