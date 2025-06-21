package routes_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/db"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
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

func TestNewUserHandler(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	assert.NotNil(t, handler)
	assert.Equal(t, mockRepo, handler.Repo)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := &routes.UserHandler{Repo: mockRepo}

	app := fiber.New()
	app.Post("/users", handler.CreateUser)

	validBody := `{
		"username": "testuser",
		"password_hash": "hashedpassword123",
		"email": "test@example.com",
		"roles": "admin"
	}`

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(validBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_FailToCreate(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := &routes.UserHandler{Repo: mockRepo}

	app := fiber.New()
	app.Post("/users", handler.CreateUser)

	validBody := `{
		"username": "testuser",
		"password_hash": "hashedpassword123",
		"email": "test@example.com",
		"roles": "admin"
	}`

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(validBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	handler := &routes.UserHandler{}

	app := fiber.New()
	app.Post("/users", handler.CreateUser)

	invalidBody := `{"username":}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := &routes.UserHandler{Repo: mockRepo}

	app := fiber.New()
	app.Get("/users", handler.GetAllUsers)

	expectedUsers := []db.ListUsersRow{
		{ID: "1", Username: "user1", Email: "user1@example.com", Roles: "admin"},
		{ID: "2", Username: "user2", Email: "user2@example.com", Roles: "user"},
	}

	mockRepo.On("ListUsers", mock.Anything).Return(expectedUsers, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var got []db.ListUsersRow
	err = json.Unmarshal(bodyBytes, &got)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, got)

	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := &routes.UserHandler{Repo: mockRepo}

	app := fiber.New()
	app.Get("/users", handler.GetAllUsers)

	mockRepo.On("ListUsers", mock.Anything).Return(nil, errors.New("db error")).Once()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
    mockRepo := new(MockUserRepo)
    handler := &routes.UserHandler{Repo: mockRepo}

    app := fiber.New()
    app.Get("/users/:id", handler.GetUserByID)

    expectedUser := db.GetUserByIDRow{
        ID:       "user-123",
        Username: "testuser",
        Email:    "test@example.com",
        Roles:    "admin",
    }

    mockRepo.On("GetUserByID", mock.Anything, "user-123").Return(expectedUser, nil).Once()

    req := httptest.NewRequest(http.MethodGet, "/users/user-123", nil)
    resp, err := app.Test(req, -1)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    bodyBytes, err := io.ReadAll(resp.Body)
    assert.NoError(t, err)
    resp.Body.Close()

    var got db.GetUserByIDRow
    err = json.Unmarshal(bodyBytes, &got)
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, got)

    mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
    mockRepo := new(MockUserRepo)
    handler := &routes.UserHandler{Repo: mockRepo}

    app := fiber.New()
    app.Get("/users/:id", handler.GetUserByID)

    mockRepo.On("GetUserByID", mock.Anything, "user-123").Return(db.GetUserByIDRow{}, errors.New("not found")).Once()

    req := httptest.NewRequest(http.MethodGet, "/users/user-123", nil)
    resp, err := app.Test(req, -1)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

    mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepo)
    handler := &routes.UserHandler{Repo: mockRepo}

    app := fiber.New()
    app.Delete("/users/:id", handler.DeleteUser)

    mockRepo.On("DeleteUser", mock.Anything, "user-123").Return(nil).Once()

    req := httptest.NewRequest(http.MethodDelete, "/users/user-123", nil)
    resp, err := app.Test(req, -1)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

    mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Fail(t *testing.T) {
    mockRepo := new(MockUserRepo)
    handler := &routes.UserHandler{Repo: mockRepo}

    app := fiber.New()
    app.Delete("/users/:id", handler.DeleteUser)

    mockRepo.On("DeleteUser", mock.Anything, "user-123").Return(errors.New("delete failed")).Once()

    req := httptest.NewRequest(http.MethodDelete, "/users/user-123", nil)
    resp, err := app.Test(req, -1)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

    mockRepo.AssertExpectations(t)
}


