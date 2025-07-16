package routes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/routes"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)
	app := fiber.New()
	app.Post("/users", handler.CreateUser)

	payload := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Roles:        "admin",
	}
	body, _ := json.Marshal(payload)

	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(p db.CreateUserParams) bool {
		return p.Username == payload.Username && p.Email == payload.Email
	})).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_InvalidBody(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)
	app := fiber.New()
	app.Post("/users", handler.CreateUser)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)
	app := fiber.New()
	app.Get("/users", handler.GetAllUsers)

	mockUsers := []db.ListUsersRow{
		{ID: uuid.New().String(), Username: "user1", Email: "user1@example.com"},
	}
	mockRepo.On("ListUsers", mock.Anything).Return(mockUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Failure(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Post("/users", handler.CreateUser)

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("db failed"))

	body := `{
		"username": "testuser",
		"email": "test@example.com",
		"passwordHash": "hashedpw",
		"roles": "admin"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	b, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(b), "failed to create user")
}

func TestGetAllUsers_Failure(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Get("/users", handler.GetAllUsers)

	mockRepo.On("ListUsers", mock.Anything).Return(nil, errors.New("db down"))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	b, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(b), "db down")
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByID)

	mockRepo.On("GetUserByID", mock.Anything, "abc").Return(db.GetUserByIDRow{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	b, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(b), "user not found")
}

func TestDeleteUser_Failure(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Delete("/users/:id", handler.DeleteUser)

	mockRepo.On("DeleteUser", mock.Anything, "abc").Return(errors.New("delete failed"))

	req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	b, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(b), "delete failed")
}

func TestCreateUser_ErrorFromRepo(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Post("/users", handler.CreateUser)

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("db error"))

	body := `{
		"username": "tester",
		"email": "tester@example.com",
		"passwordHash": "hashedpassword",
		"roles": "admin"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(respBody), "failed to create user")
}

func TestDeleteUser_ErrorFromRepo(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Delete("/users/:id", handler.DeleteUser)

	mockRepo.On("DeleteUser", mock.Anything, "u123").Return(errors.New("delete failed"))

	req := httptest.NewRequest(http.MethodDelete, "/users/u123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(respBody), "delete failed")
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByID)

	expectedUser := db.GetUserByIDRow{
		ID:       "123",
		Username: "testuser",
		Email:    "test@example.com",
		Roles:   "admin",
	}

	mockRepo.On("GetUserByID", mock.Anything, "123").Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)

	var actualUser db.GetUserByIDRow
	err = json.Unmarshal(bodyBytes, &actualUser)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	handler := routes.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Delete("/users/:id", handler.DeleteUser)

	userID := "123"
	mockRepo.On("DeleteUser", mock.Anything, userID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/users/123", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

