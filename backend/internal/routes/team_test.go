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

func mockUserMiddleware(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	}
}

func setUpTeamApp(h *routes.TeamHandler) *fiber.App {
	app := fiber.New()
	app.Use(mockUserMiddleware("user-123"))
	app.Post("/teams", h.CreateTeam)
	app.Get("/teams", h.GetTeamsByUserID)
	app.Delete("/teams", h.DeleteTeam)
	return app
}

func TestCreateTeam(t *testing.T) {

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

func TestGetTeamByUserID(t *testing.T) {

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

func TestDeleteTeam(t *testing.T) {

	mockTeamRepo := new(mocks.MockTeamRepo)
	mockWorkspaceRepo := new(mocks.MockWorkspaceRepo)
	handler := routes.NewTeamHandler(mockTeamRepo, mockWorkspaceRepo)

	teamID := "team-123"
	userID := "user-123"


	mockTeamRepo.On("GetOwnerByTeamID", mock.Anything, teamID).
		Return(userID, nil)
	mockTeamRepo.On("GetLeaderByTeamID", mock.Anything, teamID).
		Return(userID, nil)
	mockTeamRepo.On("DeleteTeam", mock.Anything, teamID).
		Return(nil)
	mockTeamRepo.On("DeleteTeamFromTeamMembers", mock.Anything, teamID).
		Return(nil)

	app := setUpTeamApp(handler)

	req := httptest.NewRequest("DELETE", "/teams/"+teamID, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "team deleted successfully")


	mockTeamRepo.AssertExpectations(t)
}
