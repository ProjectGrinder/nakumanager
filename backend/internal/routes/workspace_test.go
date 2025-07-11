package routes_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
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

func TestCreateWorkspace(t *testing.T) {
	t.Run("Create Workspace Successfully", func(t *testing.T) {
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
	})

	t.Run("Invalid request body", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/workspaces", handler.CreateWorkspace)

		payload := `{"name":`
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	})

	t.Run("Validation errors", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		tests := []struct {
			name     string
			payload  string
			expected int
		}{
			{
				name:     "Empty name",
				payload:  `{"name": ""}`,
				expected: fiber.StatusBadRequest,
			},
			{
				name:     "Missing name field",
				payload:  `{}`,
				expected: fiber.StatusBadRequest,
			},
		}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/workspaces", handler.CreateWorkspace)

		for _, tt := range tests {
			req := httptest.NewRequest("POST", "/workspaces", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			reqCtx := req.Context()
			req = req.WithContext(reqCtx)
			resp, _ := app.Test(req, -1)
			assert.Equal(t, tt.expected, resp.StatusCode, tt.name)
		}

	})

	t.Run("Create Workspace Fail", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/workspaces", handler.CreateWorkspace)

		payload := map[string]string{"name": "Test Workspace"}
		body, _ := json.Marshal(payload)

		repo.On("CreateWorkspace", mock.Anything, mock.Anything, "Test Workspace", "user-123").
			Return(errors.New("create workspace failed"))

		req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		repo.AssertCalled(t, "CreateWorkspace", mock.Anything, mock.Anything, "Test Workspace", "user-123")
	})

	t.Run("Fail to add cretor to workspace", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Post("/workspaces", handler.CreateWorkspace)

		payload := map[string]string{"name": "Test Workspace"}
		body, _ := json.Marshal(payload)

		repo.On("CreateWorkspace", mock.Anything, mock.Anything, "Test Workspace", "user-123").Return(nil)
		repo.On("AddMemberToWorkspace", mock.Anything, mock.Anything, "user-123").Return(errors.New("fail to add creator to workspace"))

		req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestGetWorkspaceByUserID(t *testing.T) {
	t.Run("Get workspace by user id successfully", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Get("/workspaces", handler.GetWorkspacesByUserID)

		expected := []db.ListWorkspacesWithMembersByUserIDRow{
			{ID: "w1", Name: "Workspace 1", OwnerID: "user-123", UserID: sql.NullString{String: "user-123", Valid: true}},
			{ID: "w2", Name: "Workspace 2", OwnerID: "user-123", UserID: sql.NullString{String: "user-123", Valid: true}},
		}
		repo.On("ListWorkspacesWithMembersByUserID", mock.Anything, "user-123").Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/workspaces", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var actual []db.ListWorkspacesWithMembersByUserIDRow
		err = json.Unmarshal(body, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Fail to get workspace by user id", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Get("/workspaces", handler.GetWorkspacesByUserID)

		repo.On("ListWorkspacesWithMembersByUserID", mock.Anything, "user-123").Return(nil, errors.New("fail to get workspace by user id"))

		req := httptest.NewRequest(http.MethodGet, "/workspaces", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestDeleteWorkspace(t *testing.T) {
	t.Run("Delete workspace successfully", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)
		repo.On("DeleteWorkspace", mock.Anything, "ws-123").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("WorkSpace ID is not provided", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

		req := httptest.NewRequest(http.MethodDelete, "/workspaces/undefined", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not found WorkSpace", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").Return(db.Workspace{}, errors.New("not found"))

		req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("User Request is not Owner of WorkSpace", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").Return(db.Workspace{ID: "ws-123", OwnerID: "user-456"}, nil)

		req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Delete WorkSpace failed", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Delete("/workspaces/:workspaceid", handler.DeleteWorkspace)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)
		repo.On("DeleteWorkspace", mock.Anything, "ws-123").Return(errors.New("failed to delete workspace"))

		req := httptest.NewRequest(http.MethodDelete, "/workspaces/ws-123", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestUpdateWorkSpace(t *testing.T) {
	t.Run("Update name workspace successfully", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)

		body := []byte(`{
		"name": "New Workspace Name"
		}`)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("RenameWorkspace", mock.Anything, "ws-123", "New Workspace Name").
			Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		respBody, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), `"message":"workspace updated successfully"`)
	})

	t.Run("Update name workspace failed", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)
		body := []byte(`{
		"name": "New Workspace Name"
		}`)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("RenameWorkspace", mock.Anything, "ws-123", "New Workspace Name").
			Return(errors.New("failed to rename workspace"))

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		respBody, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), `"error":"failed to rename workspace"`)

	})

	t.Run("Add Member to Workspace successfully", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)
		body := []byte(`{
		"add_members": ["user-456", "user-789"]
		}`)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("AddMemberToWorkspace", mock.Anything, "ws-123", "user-456").Return(nil)
		repo.On("AddMemberToWorkspace", mock.Anything, "ws-123", "user-789").Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		respBody, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), `"message":"workspace updated successfully"`)
	})

	t.Run("Failed to Add Member to Workspace successfully", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)
		body := []byte(`{
		"add_members": ["user-456", "user-789"]
		}`)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("AddMemberToWorkspace", mock.Anything, "ws-123", "user-456").Return(errors.New("failed to add member to workspace"))
		repo.On("AddMemberToWorkspace", mock.Anything, "ws-123", "user-789").Return(errors.New("failed to add member to workspace"))

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		respBody, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), `"error":"failed to add member to workspace"`)

	})

	t.Run("Remove Member from Workspace successfully", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)
		body := []byte(`{
		"remove_members": ["user-456", "user-789"]
		}`)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("RemoveMemberFromWorkspace", mock.Anything, "ws-123", "user-456").Return(nil)
		repo.On("RemoveMemberFromWorkspace", mock.Anything, "ws-123", "user-789").Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		respBody, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), `"message":"workspace updated successfully"`)
	})

	t.Run("Failed to Remove Member from Workspace successfully", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)
		body := []byte(`{
		"remove_members": ["user-456", "user-789"]
		}`)

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		repo.On("RemoveMemberFromWorkspace", mock.Anything, "ws-123", "user-456").Return(errors.New("failed to remove member from workspace"))
		repo.On("RemoveMemberFromWorkspace", mock.Anything, "ws-123", "user-789").Return(errors.New("failed to remove member from workspace"))

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		respBody, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), `"error":"failed to remove member from workspace"`)
	})

	t.Run("WorkSpace ID not provided", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)

		req := httptest.NewRequest(http.MethodPut, "/workspaces/undefined", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	})

	t.Run("Not Found WorkSpace", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", nil)
		req.Header.Set("Content-Type", "application/json")

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").Return(db.Workspace{}, errors.New("workspace not found"))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Request User Is Not Workspace Owner", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))

		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", nil)
		req.Header.Set("Content-Type", "application/json")

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").Return(db.Workspace{ID: "ws-123", OwnerID: "user-456"}, nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Invalid JSON body", func(t *testing.T) {
		repo := new(mocks.MockWorkspaceRepo)
		handler := routes.WorkspaceHandler{Repo: repo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Put("/workspaces/:workspaceid", handler.UpdateWorkspace)

		body := []byte(`{
		"name": "new name",
		}`) 

		repo.On("GetWorkspaceByID", mock.Anything, "ws-123").
			Return(db.Workspace{ID: "ws-123", OwnerID: "user-123"}, nil)

		req := httptest.NewRequest(http.MethodPut, "/workspaces/ws-123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		res, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(res), `"error":"invalid request body"`)
	})

}
