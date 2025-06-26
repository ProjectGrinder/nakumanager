package routes_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)

func setUpProjectApp(h *routes.ProjectHandler) *fiber.App {
	app := fiber.New()
	app.Use(mockUserMiddleware("user-123"))

	app.Post("/projects", h.CreateProject)
	app.Get("/projects", h.GetProjectsByUserID)
	app.Delete("/projects", h.DeleteProject)
	return app
}

func TestCreateProject(t *testing.T) {
	handler := &routes.ProjectHandler{}
	app := setUpProjectApp(handler)

	req := httptest.NewRequest("POST", "/projects", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Create Project!", string(body))
}

func TestGetProjectsByUserID(t *testing.T) {
	handler := &routes.ProjectHandler{}
	app := setUpProjectApp(handler)

	req := httptest.NewRequest("GET", "/projects", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Get Projects!", string(body))
}

func TestDeleteProject(t *testing.T) {
	handler := &routes.ProjectHandler{}
	app := setUpProjectApp(handler)

	req := httptest.NewRequest("DELETE", "/projects", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Delete Project!", string(body))
}
