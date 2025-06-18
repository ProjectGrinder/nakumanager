package routes_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)

func setUpWorkspaceApp() *fiber.App {
	app := fiber.New()
	app.Post("/workspaces", routes.CreateWorkspace)
	app.Get("/workspaces", routes.GetWorkspacesByUserID)
	app.Delete("/workspaces", routes.DeleteWorkspace)
	return app
}

func TestCreateWorkspace(t *testing.T) {
	app := setUpWorkspaceApp()
	req := httptest.NewRequest("POST", "/workspaces", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Create Workspace!", string(body))
}

func TestGetWorkspaceByUserID(t *testing.T) {
	app := setUpWorkspaceApp()
	req := httptest.NewRequest("GET", "/workspaces", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Get Workspaces!", string(body))
}

func TestDeleteWorkspace(t *testing.T) {
	app := setUpWorkspaceApp()
	req := httptest.NewRequest("DELETE", "/workspaces", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Delete Workspace!", string(body))
}
