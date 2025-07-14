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

func TestGetTeamsByUserID(t *testing.T) {
	t.Run("Get Teams Successfully", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Get("/teams", handler.GetTeamsByUserID)

		repo.On("GetTeamsByUserID", mock.Anything, "user-123").
			Return([]db.Team{{ID: "team-123", Name: "Test Team", WorkspaceID: "ws-123"}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/teams", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Get Teams Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Get("/teams", handler.GetTeamsByUserID)

		repo.On("GetTeamsByUserID", mock.Anything, "user-123").
			Return([]db.Team{}, errors.New("failed to get teams"))

		req := httptest.NewRequest(http.MethodGet, "/teams", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Team is not exists", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Get("/teams", handler.GetTeamsByUserID)

		repo.On("GetTeamsByUserID", mock.Anything, "user-123").
			Return([]db.Team{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/teams", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestDeleteTeam(t *testing.T) {
	t.Run("Delete Team Successfully", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/teams/:id", handler.DeleteTeam)

		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("user-123", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-123").Return("user-456", nil)
		repo.On("DeleteTeam", mock.Anything, "team-123").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/teams/team-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	})

	t.Run("Delete Team Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/teams/:id", handler.DeleteTeam)

		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("user-123", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-123").Return("user-456", nil)
		repo.On("DeleteTeam", mock.Anything, "team-123").Return(errors.New("failed to delete team"))

		req := httptest.NewRequest(http.MethodDelete, "/teams/team-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("TeamID is not provided", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/teams/:id", handler.DeleteTeam)

		req := httptest.NewRequest(http.MethodDelete, "/teams/empty", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Check Team Owner Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/teams/:id", handler.DeleteTeam)

		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("", errors.New("failed to check team owner"))

		req := httptest.NewRequest(http.MethodDelete, "/teams/team-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Check Team Leader Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/teams/:id", handler.DeleteTeam)

		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("user-123", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-123").Return("", errors.New("failed to check team leader"))

		req := httptest.NewRequest(http.MethodDelete, "/teams/team-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("No Permission to Remove Member", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/teams/:id", handler.DeleteTeam)

		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("user-999", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-123").Return("user-888", nil)

		req := httptest.NewRequest(http.MethodDelete, "/teams/team-123", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})
}

func TestUpdateTeam(t *testing.T) {
	t.Run("Update Team Success", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{
		"name": "Updated Team Name",
		"addMembers": ["user-new"],
		"removeMembers": ["user-old"],
		"newLeaderID": "user-new"
	}`

		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-xyz").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-xyz").Return("user-leader", nil)
		repo.On("RenameTeam", mock.Anything, db.RenameTeamParams{ID: "team-xyz", Name: "Updated Team Name"}).Return(nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-xyz", "user-new").Return(false, nil)
		repo.On("AddMemberToTeam", mock.Anything, db.AddMemberToTeamParams{TeamID: "team-xyz", UserID: "user-new"}).Return(nil)
		repo.On("RemoveMemberFromTeam", mock.Anything, db.RemoveMemberFromTeamParams{TeamID: "team-xyz", UserID: "user-old"}).Return(nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-xyz", "user-new").Return(true, nil)
		repo.On("SetLeaderToTeam", mock.Anything, db.SetLeaderToTeamParams{ID: "team-xyz", LeaderID: "user-new"}).Return(nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("UpdateTeam Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{
		"name": "Updated Team Name"
		}`

		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-xyz").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-xyz").Return("user-leader", nil)
		repo.On("RenameTeam", mock.Anything, db.RenameTeamParams{ID: "team-xyz", Name: "Updated Team Name"}).Return(errors.New("failed to rename team"))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		repo.AssertExpectations(t)
	})

	t.Run("TeamID is not provided", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		req := httptest.NewRequest(http.MethodPatch, "/teams/empty", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Bad Body Request", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		req := httptest.NewRequest(http.MethodPatch, "/teams/team-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	})

	t.Run("Check Team Owner Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		repo.On("IsTeamExists", mock.Anything, "team-123").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("", errors.New("failed to check team owner"))

		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		repo.AssertExpectations(t)
	})

	t.Run("Check Team Leader Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		repo.On("IsTeamExists", mock.Anything, "team-123").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("user-123", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-123").Return("", errors.New("failed to check team leader"))

		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("No Permission to Remove Member", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		repo.On("IsTeamExists", mock.Anything, "team-123").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("user-999", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-123").Return("user-888", nil)

		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Check Team Exists Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-123").Return(false, errors.New("DB error"))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		repo.AssertExpectations(t)
	})

	t.Run("Team Does Not Exist", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"name": "x"}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(false, nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		repo.AssertExpectations(t)
	})

	t.Run("AddMemberToTeam - Success", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"add_members": ["user-new"]}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-xyz").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-xyz").Return("user-leader", nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-xyz", "user-new").Return(false, nil)
		repo.On("AddMemberToTeam", mock.Anything, db.AddMemberToTeamParams{
			TeamID: "team-xyz", UserID: "user-new",
		}).Return(nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		repo.AssertExpectations(t)
	})

	t.Run("AddMember - Already Exists (Should Skip)", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"add_members": ["user-existing"]}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-xyz").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-xyz").Return("user-leader", nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-xyz", "user-existing").Return(true, nil)

		repo.On("AddMemberToTeam", mock.Anything, mock.Anything).Maybe() 

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		repo.AssertNotCalled(t, "AddMemberToTeam", mock.Anything, mock.Anything)
	})

	t.Run("AddMemberToTeam - Insert Failed", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"add_members": ["user-new"]}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-xyz").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-xyz").Return("user-leader", nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-xyz", "user-new").Return(false, nil)
		repo.On("AddMemberToTeam", mock.Anything, db.AddMemberToTeamParams{
			TeamID: "team-xyz", UserID: "user-new",
		}).Return(errors.New("insert failed"))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		repo.AssertExpectations(t)
	})

	t.Run("AddMembers - IsMemberInTeam Error", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"add_members": ["user-new"]}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-123").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-123").Return("user-leader", nil)

		repo.On("IsMemberInTeam", mock.Anything, "team-123", "user-new").Return(false, errors.New("DB error"))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("RemoveMembers - Success", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"remove_members": ["user-old"]}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-abc", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-abc").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-abc").Return("user-123", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-abc").Return("user-leader", nil)

		repo.On("RemoveMemberFromTeam", mock.Anything, db.RemoveMemberFromTeamParams{
			TeamID: "team-abc", UserID: "user-old",
		}).Return(nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		repo.AssertExpectations(t)
	})

	t.Run("RemoveMembers - RemoveMember Error", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"remove_members": ["user-old"]}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-abc", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-abc").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-abc").Return("user-123", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-abc").Return("user-leader", nil)

		repo.On("RemoveMemberFromTeam", mock.Anything, db.RemoveMemberFromTeamParams{
			TeamID: "team-abc", UserID: "user-old",
		}).Return(errors.New("remove failed"))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		repo.AssertExpectations(t)
	})

	t.Run("Set New Leader Success", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"new_leader_id": "user-new"}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-xyz").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-xyz").Return("user-leader", nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-xyz", "user-new").Return(true, nil)
		repo.On("SetLeaderToTeam", mock.Anything, db.SetLeaderToTeamParams{
			ID:       "team-xyz",
			LeaderID: "user-new",
		}).Return(nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		repo.AssertExpectations(t)
	})

	t.Run("Set New Leader - Check Member Error", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"new_leader_id": "user-new"}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-xyz").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-xyz").Return("user-leader", nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-xyz", "user-new").
			Return(false, errors.New("db error"))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Set New Leader - Not a Member", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"new_leader_id": "user-new"}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-xyz", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		repo.On("IsTeamExists", mock.Anything, "team-xyz").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-xyz").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-xyz").Return("user-leader", nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-xyz", "user-new").
			Return(false, nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("SetLeaderToTeam - Failed to Save", func(t *testing.T) {
		repo := new(mocks.MockTeamRepository)
		workSpaceRepo := new(mocks.MockWorkspaceRepo)
		handler := routes.TeamHandler{Repo: repo, WorkspaceRepo: workSpaceRepo}

		app := fiber.New()
		app.Use(withUserID("user-owner"))
		app.Patch("/teams/:id", handler.UpdateTeam)

		reqBody := `{"new_leader_id": "user-new"}`
		req := httptest.NewRequest(http.MethodPatch, "/teams/team-123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		repo.On("IsTeamExists", mock.Anything, "team-123").Return(true, nil)
		repo.On("GetOwnerByTeamID", mock.Anything, "team-123").Return("user-owner", nil)
		repo.On("GetLeaderByTeamID", mock.Anything, "team-123").Return("user-leader", nil)
		repo.On("IsMemberInTeam", mock.Anything, "team-123", "user-new").Return(true, nil)

		repo.On("SetLeaderToTeam", mock.Anything, db.SetLeaderToTeamParams{
			ID:       "team-123",
			LeaderID: "user-new",
		}).Return(errors.New("mock DB error"))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		repo.AssertExpectations(t)
	})
}
