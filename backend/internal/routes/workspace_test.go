package routes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/db"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nack098/nakumanager/internal/routes"
)

func withUserID(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	}
}

func TestNewWorkspaceHandler(t *testing.T) {
	mockWorkspaceRepo := new(mocks.MockWorkspaceRepo)
	mockUserRepo := new(mocks.MockUserRepo)

	handler := routes.NewWorkspaceHandler(mockWorkspaceRepo, mockUserRepo)

	assert.NotNil(t, handler)
	assert.Equal(t, mockWorkspaceRepo, handler.Repo)
	assert.Equal(t, mockUserRepo, handler.UserRepo)
}

func TestCreateWorkspace_Success(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces", handler.CreateWorkspace)

	payload := map[string]string{"name": "Test Workspace"}
	body, _ := json.Marshal(payload)

	repo.On("CreateWorkspace", mock.Anything, mock.Anything, "Test Workspace", "user-123").Return(nil)
	repo.On("AddMemberToWorkspace", mock.Anything, mock.Anything, "user-123").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestCreateWorkspace_InvalidJSON(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := routes.WorkspaceHandler{Repo: repo}
	app := fiber.New()
	app.Post("/workspaces", handler.CreateWorkspace)

	req := httptest.NewRequest(http.MethodPost, "/workspaces", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateWorkspace_MissingUserID(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Post("/workspaces", handler.CreateWorkspace)

	payload := map[string]string{"name": "Test Workspace"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestCreateWorkspace_ValidationError(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces", handler.CreateWorkspace)

	payload := map[string]string{"name": ""} 
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateWorkspace_CreateWorkspaceFail(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces", handler.CreateWorkspace)

	payload := map[string]string{"name": "Test Workspace"}
	body, _ := json.Marshal(payload)

	repo.On("CreateWorkspace", mock.Anything, mock.Anything, "Test Workspace", "user-123").Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestCreateWorkspace_AddMemberFail(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces", handler.CreateWorkspace)

	payload := map[string]string{"name": "Test Workspace"}
	body, _ := json.Marshal(payload)

	repo.On("CreateWorkspace", mock.Anything, mock.Anything, "Test Workspace", "user-123").Return(nil)
	repo.On("AddMemberToWorkspace", mock.Anything, mock.Anything, "user-123").Return(errors.New("add member error"))

	req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestGetWorkspacesByUserID_Success(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	expected := []db.ListWorkspacesWithMembersByUserIDRow{
		{ID: "w1", Name: "Workspace 1"},
		{ID: "w2", Name: "Workspace 2"},
	}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Get("/workspaces", handler.GetWorkspacesByUserID)

	repo.On("ListWorkspacesWithMembersByUserID", mock.Anything, "user-123").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/workspaces", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetWorkspacesByUserID_MissingUserID(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Get("/workspaces", handler.GetWorkspacesByUserID)

	req := httptest.NewRequest(http.MethodGet, "/workspaces", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestGetWorkspacesByUserID_RepoError(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Get("/workspaces", handler.GetWorkspacesByUserID)

	repo.On("ListWorkspacesWithMembersByUserID", mock.Anything, "user-123").Return(nil, errors.New("db fail"))

	req := httptest.NewRequest(http.MethodGet, "/workspaces", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteWorkspace_WorkspaceIDRequired(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Delete("/workspaces/:workspaceid?", handler.DeleteWorkspace)

	req := httptest.NewRequest(http.MethodDelete, "/workspaces/", nil) 
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDeleteWorkspace_Unauthorized(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

	req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestDeleteWorkspace_WorkspaceNotFound(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
		Return(db.Workspace{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestDeleteWorkspace_Forbidden(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

	ws := db.Workspace{ID: "ws-123", OwnerID: "user-999"}
	repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
		Return(ws, nil)

	req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestDeleteWorkspace_DeleteFailed(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

	ws := db.Workspace{ID: "ws-123", OwnerID: "user-123"}
	repo.On("GetWorkspaceByID", mock.Anything, "ws-123").Return(ws, nil)
	repo.On("DeleteWorkspace", mock.Anything, "ws-123").Return(errors.New("delete failed"))

	req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteWorkspace_Success(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

	ws := db.Workspace{ID: "ws-123", OwnerID: "user-123"}
	repo.On("GetWorkspaceByID", mock.Anything, "ws-123").Return(ws, nil)
	repo.On("DeleteWorkspace", mock.Anything, "ws-123").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAddMemberToWorkspace_WorkspaceIDRequired(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members", handler.AddMemberToWorkspace)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/", nil) 
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode) 
}

func TestAddMemberToWorkspace_Unauthorized(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Post("/workspaces/:workspaceid/members", handler.AddMemberToWorkspace)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAddMemberToWorkspace_WorkspaceNotFound(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members", handler.AddMemberToWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(db.Workspace{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestAddMemberToWorkspace_Forbidden(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-999"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members", handler.AddMemberToWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestAddMemberToWorkspace_BadRequest_EmptyUserID(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members", handler.AddMemberToWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)

	// Missing user_id in body
	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAddMemberToWorkspace_AddMemberFail(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members", handler.AddMemberToWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)
	repo.On("AddMemberToWorkspace", mock.Anything, "ws-1", "user-456").Return(errors.New("fail add"))

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestAddMemberToWorkspace_Success(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members", handler.AddMemberToWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)
	repo.On("AddMemberToWorkspace", mock.Anything, "ws-1", "user-456").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAddMemberToWorkspace_WorkspaceIDEmpty(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members", handler.AddMemberToWorkspace)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/empty/members", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRemoveMemberFromWorkspace_WorkspaceIDRequired(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members/remove", handler.RemoveMemberFromWorkspace)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/empty/members/remove", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRemoveMemberFromWorkspace_Unauthorized(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Post("/workspaces/:workspaceid/members/remove", handler.RemoveMemberFromWorkspace)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members/remove", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestRemoveMemberFromWorkspace_WorkspaceNotFound(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members/remove", handler.RemoveMemberFromWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(db.Workspace{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members/remove", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestRemoveMemberFromWorkspace_Forbidden(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-999"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members/remove", handler.RemoveMemberFromWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members/remove", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRemoveMemberFromWorkspace_BadRequest_EmptyUserID(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members/remove", handler.RemoveMemberFromWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members/remove", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRemoveMemberFromWorkspace_RemoveMemberFail(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members/remove", handler.RemoveMemberFromWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)
	repo.On("RemoveMemberFromWorkspace", mock.Anything, "ws-1", "user-456").Return(errors.New("fail remove"))

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members/remove", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestRemoveMemberFromWorkspace_Success(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/workspaces/:workspaceid/members/remove", handler.RemoveMemberFromWorkspace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)
	repo.On("RemoveMemberFromWorkspace", mock.Anything, "ws-1", "user-456").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/members/remove", strings.NewReader(`{"user_id":"user-456"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRenameWorkSpace_WorkspaceIDRequired(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Put("/workspaces/:workspaceid/rename", handler.RenameWorkSpace)

	req := httptest.NewRequest(http.MethodPut, "/workspaces/empty/rename", strings.NewReader(`{"name":"New Name"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRenameWorkSpace_Unauthorized(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Put("/workspaces/:workspaceid/rename", handler.RenameWorkSpace)

	req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-1/rename", strings.NewReader(`{"name":"New Name"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestRenameWorkSpace_WorkspaceNotFound(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Put("/workspaces/:workspaceid/rename", handler.RenameWorkSpace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(db.Workspace{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-1/rename", strings.NewReader(`{"name":"New Name"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestRenameWorkSpace_Forbidden(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-999"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Put("/workspaces/:workspaceid/rename", handler.RenameWorkSpace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)

	req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-1/rename", strings.NewReader(`{"name":"New Name"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestRenameWorkSpace_BadRequest_NameRequired(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Put("/workspaces/:workspaceid/rename", handler.RenameWorkSpace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)

	req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-1/rename", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRenameWorkSpace_RenameFail(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Put("/workspaces/:workspaceid/rename", handler.RenameWorkSpace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)
	repo.On("RenameWorkspace", mock.Anything, "ws-1", "New Name").Return(errors.New("fail rename"))

	req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-1/rename", strings.NewReader(`{"name":"New Name"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestRenameWorkSpace_Success(t *testing.T) {
	repo := new(mocks.MockWorkspaceRepo)
	handler := &routes.WorkspaceHandler{Repo: repo}

	ws := db.Workspace{ID: "ws-1", OwnerID: "user-123"}

	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Put("/workspaces/:workspaceid/rename", handler.RenameWorkSpace)

	repo.On("GetWorkspaceByID", mock.Anything, "ws-1").Return(ws, nil)
	repo.On("RenameWorkspace", mock.Anything, "ws-1", "New Name").Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-1/rename", strings.NewReader(`{"name":"New Name"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
