package routes_test

import (
	"net/http/httptest"
	"testing"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)


func setUpProjectApp() *fiber.App {
	app := fiber.New()
	app.Post("/projects", routes.CreateProject)
	app.Get("/projects", routes.GetProjectsByUserID)
	app.Delete("/projects", routes.DeleteProject)
	return app
}


func TestCreateProject(t *testing.T) {
	app := setUpProjectApp()
	req := httptest.NewRequest("POST", "/projects", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Create Project!", string(body))
}


func TestGetProjectsByUserID(t *testing.T) {
	app := setUpProjectApp()
	req := httptest.NewRequest("GET", "/projects", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Get Projects!", string(body))
}

func TestDeleteProject(t *testing.T) {
	app := setUpProjectApp()
	req := httptest.NewRequest("DELETE", "/projects", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Delete Project!", string(body))
}
