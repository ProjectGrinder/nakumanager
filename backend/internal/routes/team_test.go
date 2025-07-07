package routes_test

import (
	"io"
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

func setUpTeamApp(h *routes.TeamHandler) *fiber.App {
	app := fiber.New()
	app.Use(withUserID("user-123"))
	app.Post("/teams", h.CreateTeam)
	app.Get("/teams", h.GetTeamsByUserID)
	app.Post("/teams/:id/members", h.AddMemberToTeam)
	app.Delete("/teams/:id/members", h.RemoveMemberFromTeam)
	app.Post("/teams/:id/rename", h.RenameTeam)
	app.Post("/teams/:id/leader", h.SetTeamLeader)
	app.Delete("/teams/:id", h.DeleteTeam)
	return app
}

func TestCreateTeam_Success(t *testing.T) {

	mockTeamRepo := new(mocks.MockTeamRepo)
	mockWorkspaceRepo := new(mocks.MockWorkspaceRepo)

	handler := routes.NewTeamHandler(mockTeamRepo, mockWorkspaceRepo)

	mockWorkspaceRepo.On("GetWorkspaceByID", mock.Anything, "ws-123").
		Return(db.Workspace{
			ID:      "ws-123",
			OwnerID: "user-123",
		}, nil)

	mockTeamRepo.On("CreateTeam", mock.Anything, mock.Anything).
		Return(nil)

	mockTeamRepo.On("AddMemberToTeam", mock.Anything, mock.Anything).
		Return(nil)

	app := setUpTeamApp(handler)

	body := `{"name":"Team A", "workspace_id":"ws-123"}`
	req := httptest.NewRequest("POST", "/teams", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(respBody), "team created successfully")

	mockWorkspaceRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestGetTeamByUserID_Success(t *testing.T) {

	mockTeamRepo := new(mocks.MockTeamRepo)
	mockWorkspaceRepo := new(mocks.MockWorkspaceRepo)
	handler := routes.NewTeamHandler(mockTeamRepo, mockWorkspaceRepo)

	mockTeamRepo.On("GetTeamsByUserID", mock.Anything, "user-123").
		Return([]db.Team{
			{ID: "team-1", Name: "Team A"},
		}, nil)

	app := setUpTeamApp(handler)

	req := httptest.NewRequest("GET", "/teams", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"Team A"`)
}

func TestDeleteTeam_Success(t *testing.T) {
	mockTeamRepo := new(mocks.MockTeamRepo)
	mockWorkspaceRepo := new(mocks.MockWorkspaceRepo)
	handler := routes.NewTeamHandler(mockTeamRepo, mockWorkspaceRepo)

	teamID := "team-123"
	userID := "user-123"

	mockTeamRepo.On("GetOwnerByTeamID", mock.Anything, teamID).Return(userID, nil)
	mockTeamRepo.On("GetLeaderByTeamID", mock.Anything, teamID).Return(userID, nil)
	mockTeamRepo.On("DeleteTeam", mock.Anything, teamID).Return(nil)

	app := setUpTeamApp(handler)

	req := httptest.NewRequest("DELETE", "/teams/"+teamID, nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "team deleted successfully")

	mockTeamRepo.AssertExpectations(t)
}

func TestAddMemberToTeam_Success(t *testing.T) {
	mockTeamRepo := new(mocks.MockTeamRepo)
	mockWorkspaceRepo := new(mocks.MockWorkspaceRepo)
	handler := routes.NewTeamHandler(mockTeamRepo, mockWorkspaceRepo)

	teamID := "team-123"
	userID := "user-456"

	mockTeamRepo.On("AddMemberToTeam", mock.Anything, teamID, userID).Return(nil)

	app := setUpTeamApp(handler)

	req := httptest.NewRequest("POST", "/teams/"+teamID+"/members", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "member added to team successfully")

	mockTeamRepo.AssertExpectations(t)
}
