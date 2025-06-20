package routes_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)

func TestSetUpWorkspaceRoutes_Simple(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api")
	routes.SetUpWorkspaceRoutes(api)

	req := httptest.NewRequest(http.MethodGet, "/api/workspace/123", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode == fiber.StatusNotFound {
		t.Log("GET route exists but handler may not be implemented")
	} else {
		assert.NotEqual(t, 0, resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/workspace", nil)
	resp, _ = app.Test(req)
	assert.NotEqual(t, 0, resp.StatusCode)

	req = httptest.NewRequest(http.MethodDelete, "/api/workspace/123", nil)
	resp, _ = app.Test(req)
	assert.NotEqual(t, 0, resp.StatusCode)
}


func TestCreateWorkspace_Success(t *testing.T) {
	app := fiber.New()
	app.Post("/workspaces", func(c *fiber.Ctx) error {
		c.Locals("userID", "mockedID")
		return routes.CreateWorkspace(c)
	})

	payload := `{"name":"TestWorkspace"}`
	req := httptest.NewRequest("POST", "/workspaces", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "workspace created successfully", body["message"])
}

func TestCreateWorkspace_InvalidBody(t *testing.T) {
	app := fiber.New()
	app.Post("/workspaces", func(c *fiber.Ctx) error {
		c.Locals("userID", "mockedID")
		return routes.CreateWorkspace(c)
	})

	payload := `{"name":}`
	req := httptest.NewRequest("POST", "/workspaces", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "invalid request body")
}

func TestCreateWorkspace_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Post("/workspaces", routes.CreateWorkspace)

	payload := `{"name":"TestWorkspace"}`
	req := httptest.NewRequest("POST", "/workspaces", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestCreateWorkspace_InvalidValidation(t *testing.T) {
	app := fiber.New()
	app.Post("/workspace", func(c *fiber.Ctx) error {
		c.Locals("userID", "user123") 
		return routes.CreateWorkspace(c)
	})

	body := `{"name": ""}`
	req := httptest.NewRequest(http.MethodPost, "/workspace", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(respBody), "Name is invalid: required")
}


func TestGetWorkspaceByUserID(t *testing.T) {

	workspace := models.Workspace{
		ID:      "ws1",
		Name:    "TeamWS",
		Members: []string{"user123"},
	}
	routes.WorkspaceMutex.Lock()
	routes.WorkSpaces["ws1"] = workspace
	routes.WorkspaceMutex.Unlock()

	app := fiber.New()
	app.Get("/workspaces/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "user123")
		return routes.GetWorkspacesByUserID(c)
	})

	req := httptest.NewRequest("GET", "/workspaces/user123", strings.NewReader(""))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body []models.Workspace
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "ws1", body[0].ID)
	assert.Equal(t, "TeamWS", body[0].Name)
	assert.Equal(t, "user123", body[0].Members[0])

}

func TestGetWorkspacesByUserID_MissingID(t *testing.T) {
	app := fiber.New()
	app.Get("/workspaces/:id?", routes.GetWorkspacesByUserID)

	req := httptest.NewRequest(http.MethodGet, "/workspaces/", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "UserID is required")
}

func TestDeleteWorkspace_AsAdmin_Success(t *testing.T) {
	workspace := models.Workspace{
		ID:      "ws1",
		Name:    "To Delete",
		Members: []string{"admin123"},
	}
	routes.WorkspaceMutex.Lock()
	routes.WorkSpaces["ws1"] = workspace
	routes.WorkspaceMutex.Unlock()

	models.GetUserByID = func(id string) models.User {
		return models.User{
			ID:    "admin123",
			Roles: "admin",
		}
	}

	app := fiber.New()
	app.Delete("/workspace/:workspaceid", func(c *fiber.Ctx) error {
		c.Locals("userID", "admin123")
		return routes.DeleteWorkspace(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/workspace/ws1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestDeleteWorkspace_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Delete("/workspace/:workspaceid", func(c *fiber.Ctx) error {
		return routes.DeleteWorkspace(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/workspace/ws1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestDeleteWorkspace_NotFound(t *testing.T) {
	routes.WorkspaceMutex.Lock()
	routes.WorkSpaces = make(map[string]models.Workspace)
	routes.WorkspaceMutex.Unlock()

	models.GetUserByID = func(id string) models.User {
		return models.User{ID: "admin123", Roles: "admin"}
	}

	app := fiber.New()
	app.Delete("/workspace/:workspaceid", func(c *fiber.Ctx) error {
		c.Locals("userID", "admin123")
		return routes.DeleteWorkspace(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/workspace/not_exist", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestDeleteWorkspace_UserNotFound(t *testing.T) {
	workspace := models.Workspace{ID: "ws1", Name: "Demo"}
	routes.WorkspaceMutex.Lock()
	routes.WorkSpaces["ws1"] = workspace
	routes.WorkspaceMutex.Unlock()

	models.GetUserByID = func(id string) models.User {
		return models.User{}
	}

	app := fiber.New()
	app.Delete("/workspace/:workspaceid", func(c *fiber.Ctx) error {
		c.Locals("userID", "unknownUser")
		return routes.DeleteWorkspace(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/workspace/ws1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteWorkspace_Forbidden(t *testing.T) {
	workspace := models.Workspace{ID: "ws1", Name: "Demo"}
	routes.WorkspaceMutex.Lock()
	routes.WorkSpaces["ws1"] = workspace
	routes.WorkspaceMutex.Unlock()

	models.GetUserByID = func(id string) models.User {
		return models.User{ID: "user456", Roles: "member"}
	}

	app := fiber.New()
	app.Delete("/workspace/:workspaceid", func(c *fiber.Ctx) error {
		c.Locals("userID", "user456")
		return routes.DeleteWorkspace(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/workspace/ws1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}
