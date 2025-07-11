package routes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nack098/nakumanager/internal/db"
	"github.com/nack098/nakumanager/internal/routes"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
)

func TestNewTeamHandler(t *testing.T) {
	mockTeamRepo := new(mocks.MockTeamRepository)
	mockWorkspaceRepo := new(mocks.MockWorkspaceRepo)
	handler := routes.NewTeamHandler(mockTeamRepo, mockWorkspaceRepo)

	assert.NotNil(t, handler)
	assert.Equal(t, mockTeamRepo, handler.Repo)
	assert.Equal(t, mockWorkspaceRepo, handler.WorkspaceRepo)

}

func TestCreateTeam(t *testing.T) {
	t.Run("Create Team Successfully", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/teams", handler.CreateTeam)

		payload := map[string]string{"name": "Test Team", "workspace_id": "ws-123"}
		body, _ := json.Marshal(payload)

		workSpaceRepo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("CreateTeam", mock.Anything, mock.Anything).Return(nil)
		repo.On("AddMemberToTeam", mock.Anything, mock.Anything).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Fail to create team", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/teams", handler.CreateTeam)

		payload := map[string]string{"name": "Test Team", "workspace_id": "ws-123"}
		body, _ := json.Marshal(payload)

		workSpaceRepo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("CreateTeam", mock.Anything, mock.Anything).Return(errors.New("failed to add member to team"))

		req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Fail to add member to team", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/teams", handler.CreateTeam)

		payload := map[string]string{"name": "Test Team", "workspace_id": "ws-123"}
		body, _ := json.Marshal(payload)

		workSpaceRepo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("CreateTeam", mock.Anything, mock.Anything).Return(nil)
		repo.On("AddMemberToTeam", mock.Anything, mock.Anything).Return(errors.New("failed to add member to team"))

		req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/teams", handler.CreateTeam)

		body := []byte(`{
		"name": "Test Team",
		}`)

		workSpaceRepo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Workspace Not Found", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/teams", handler.CreateTeam)

		payload := map[string]string{"name": "Test Team", "workspace_id": "ws-123"}
		body, _ := json.Marshal(payload)

		workSpaceRepo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{}, errors.New("workspace not found"))

		req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("User is not workspace owner", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/teams", handler.CreateTeam)

		payload := map[string]string{"name": "Test Team", "workspace_id": "ws-123"}
		body, _ := json.Marshal(payload)

		workSpaceRepo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-456"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Create team validation fails", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/teams", handler.CreateTeam)

		payload := map[string]string{"name": "", "workspace_id": "ws-123"}
		body, _ := json.Marshal(payload)

		workSpaceRepo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

}
