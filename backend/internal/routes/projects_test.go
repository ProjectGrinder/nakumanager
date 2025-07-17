package routes_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/routes"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
)

func TestNewProjectHandler(t *testing.T) {
	db := &sql.DB{}
	projectRepo := new(mocks.MockProjectRepo)
	teamRepo := new(mocks.MockTeamRepository)

	handler := routes.NewProjectHandler(db, projectRepo, teamRepo)

	assert.Equal(t, db, handler.DB)
	assert.Equal(t, projectRepo, handler.Repo)
	assert.Equal(t, teamRepo, handler.TeamRepo)
}

func TestCreateProject(t *testing.T) {
	app := fiber.New()

	mockRepo := new(mocks.MockProjectRepo)
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := &routes.ProjectHandler{Repo: mockRepo, TeamRepo: mockTeamRepo}

	app.Post("/projects", func(c *fiber.Ctx) error {
		c.Locals("userID", "user-123")
		return handler.CreateProject(c)
	})

	t.Run("Create Project Successfully", func(t *testing.T) {
		body := models.CreateProject{
			TeamID:      "team-1",
			WorkspaceID: "workspace-1",
			Name:        "Test Project",
		}
		jsonBody, _ := json.Marshal(body)

		team := db.Team{ID: "team-1", WorkspaceID: "workspace-1"}
		mockTeamRepo.On("GetTeamByID", mock.Anything, "team-1").Return(team, nil)
		mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(true, nil)
		mockRepo.On("CreateProject", mock.Anything, mock.Anything).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Create Project Failed", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockTeamRepo := new(mocks.MockTeamRepository)
		handler := &routes.ProjectHandler{Repo: mockRepo, TeamRepo: mockTeamRepo}

		app.Post("/projects", func(c *fiber.Ctx) error {
			c.Locals("userID", "user-123")
			return handler.CreateProject(c)
		})

		body := models.CreateProject{
			TeamID:      "team-1",
			WorkspaceID: "workspace-1",
			Name:        "Test Project",
		}
		jsonBody, _ := json.Marshal(body)

		mockTeamRepo.On("GetTeamByID", mock.Anything, "team-1").Return(db.Team{ID: "team-1", WorkspaceID: "workspace-1"}, nil)
		mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(true, nil)
		mockRepo.On("CreateProject", mock.Anything, mock.Anything).Return(errors.New("failed to create project"))

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "failed to create project", response["error"])
	})

	t.Run("Invalid request body", func(t *testing.T) {
		reqBody := `{"name":}`
		req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("team not found", func(t *testing.T) {
		mockTeamRepo := new(mocks.MockTeamRepository)
		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{
			Repo:     mockRepo,
			TeamRepo: mockTeamRepo,
		}

		app := fiber.New()
		app.Post("/projects", func(c *fiber.Ctx) error {
			c.Locals("userID", "user-123")
			return handler.CreateProject(c)
		})

		body := models.CreateProject{
			TeamID:      "team-1",
			WorkspaceID: "workspace-1",
		}
		jsonBody, _ := json.Marshal(body)

		mockTeamRepo.On("GetTeamByID", mock.Anything, "team-1").
			Return(db.Team{}, sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("workspace mismatch", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockTeamRepo := new(mocks.MockTeamRepository)
		handler := &routes.ProjectHandler{Repo: mockRepo, TeamRepo: mockTeamRepo}

		app.Post("/projects", func(c *fiber.Ctx) error {
			c.Locals("userID", "user-123")
			return handler.CreateProject(c)
		})

		mockTeamRepo.On("GetTeamByID", mock.Anything, "team-1").
			Return(db.Team{ID: "team-1", WorkspaceID: "workspace-XYZ"}, nil)

		body := models.CreateProject{
			TeamID:      "team-1",
			WorkspaceID: "workspace-ABC",
			Name:        "Mismatch Project",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "team does not belong to the specified workspace", response["error"])
	})

	t.Run("failed to check team membership", func(t *testing.T) {
		mockTeamRepo := new(mocks.MockTeamRepository)
		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{
			Repo:     mockRepo,
			TeamRepo: mockTeamRepo,
		}

		app := fiber.New()
		app.Post("/projects", func(c *fiber.Ctx) error {
			c.Locals("userID", "user-123")
			return handler.CreateProject(c)
		})

		body := models.CreateProject{
			TeamID:      "team-1",
			WorkspaceID: "workspace-1",
		}
		jsonBody, _ := json.Marshal(body)

		mockTeamRepo.On("GetTeamByID", mock.Anything, "team-1").
			Return(db.Team{ID: "team-1", WorkspaceID: "workspace-1"}, nil)

		mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").
			Return(false, errors.New("failed to check team membership"))

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("User not a member of team", func(t *testing.T) {
		mockTeamRepo := new(mocks.MockTeamRepository)
		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{
			Repo:     mockRepo,
			TeamRepo: mockTeamRepo,
		}

		app := fiber.New()
		app.Post("/projects", func(c *fiber.Ctx) error {
			c.Locals("userID", "user-123")
			return handler.CreateProject(c)
		})

		body := models.CreateProject{
			TeamID:      "team-1",
			WorkspaceID: "workspace-1",
		}
		jsonBody, _ := json.Marshal(body)

		mockTeamRepo.On("GetTeamByID", mock.Anything, "team-1").
			Return(db.Team{ID: "team-1", WorkspaceID: "workspace-1"}, nil)

		mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").
			Return(false, nil)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("empty LeaderID should become nil", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockTeamRepo := new(mocks.MockTeamRepository)
		handler := &routes.ProjectHandler{Repo: mockRepo, TeamRepo: mockTeamRepo}

		app.Post("/projects", func(c *fiber.Ctx) error {
			c.Locals("userID", "user-123")
			return handler.CreateProject(c)
		})

		mockTeamRepo.On("GetTeamByID", mock.AnythingOfType("*fasthttp.RequestCtx"), "team-1").
			Return(db.Team{ID: "team-1", WorkspaceID: "workspace-1"}, nil)

		mockTeamRepo.On("IsMemberInTeam", mock.AnythingOfType("*fasthttp.RequestCtx"), "team-1", "user-123").
			Return(true, nil)

		leaderID := ""
		body := models.CreateProject{
			TeamID:      "team-1",
			WorkspaceID: "workspace-1",
			Name:        "Project with empty Leader",
			LeaderID:    &leaderID,
		}

		jsonBody, _ := json.Marshal(body)

		mockRepo.On("CreateProject", mock.Anything, mock.MatchedBy(func(p models.CreateProject) bool {
			return p.LeaderID == nil
		})).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

}

func TestGetProjectByUserID(t *testing.T) {
	t.Run("Get Project Successfully", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockTeamRepo := new(mocks.MockTeamRepository)
		handler := &routes.ProjectHandler{Repo: mockRepo, TeamRepo: mockTeamRepo}

		app.Use(withUserID("user-123"))
		app.Get("/projects", handler.GetProjectsByUserID)

		mockRepo.On("GetProjectsByUserID", mock.Anything, "user-123").Return([]models.CreateProject{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Get Project Failed", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockTeamRepo := new(mocks.MockTeamRepository)
		handler := &routes.ProjectHandler{Repo: mockRepo, TeamRepo: mockTeamRepo}

		app.Use(withUserID("user-123"))
		app.Get("/projects", handler.GetProjectsByUserID)

		mockRepo.On("GetProjectsByUserID", mock.Anything, "user-123").Return("", errors.New("failed to get projects"))

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestDeleteProject(t *testing.T) {
	t.Run("Delete Project Successfully", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockTeamRepo := new(mocks.MockTeamRepository)
		handler := &routes.ProjectHandler{Repo: mockRepo, TeamRepo: mockTeamRepo}

		app.Use(withUserID("user-123"))
		app.Delete("/projects/:id", handler.DeleteProject)

		mockRepo.On("GetProjectByID", mock.Anything, "project-123").
			Return(db.Project{ID: "project-123", CreatedBy: "user-123"}, nil)

		mockRepo.On("DeleteProject", mock.Anything, "project-123").
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/projects/project-123", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Project deleted successfully", response["message"])
	})

	t.Run("Delete Project Failed", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Delete("/projects/:id", handler.DeleteProject)

		mockRepo.On("GetProjectByID", mock.Anything, "project-err").
			Return(db.Project{ID: "project-err", CreatedBy: "user-123"}, nil)

		mockRepo.On("DeleteProject", mock.Anything, "project-err").
			Return(errors.New("delete failed"))

		req := httptest.NewRequest(http.MethodDelete, "/projects/project-err", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "failed to delete project", response["error"])
	})

	t.Run("Project Id is not provided", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Delete("/projects/:id", handler.DeleteProject)

		req := httptest.NewRequest(http.MethodDelete, "/projects/undefined", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "project_id is required", response["error"])
	})

	t.Run("Not found project", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Delete("/projects/:id", handler.DeleteProject)

		mockRepo.On("GetProjectByID", mock.Anything, "project-123").
			Return(db.Project{}, errors.New("project not found"))

		req := httptest.NewRequest(http.MethodDelete, "/projects/project-123", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "failed to fetch project", response["error"])
	})

	t.Run("User is not Owner", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockTeamRepo := new(mocks.MockTeamRepository)
		handler := &routes.ProjectHandler{Repo: mockRepo, TeamRepo: mockTeamRepo}

		app.Use(withUserID("user-123"))
		app.Delete("/projects/:id", handler.DeleteProject)

		mockRepo.On("GetProjectByID", mock.Anything, "project-123").
			Return(db.Project{ID: "project-123", CreatedBy: "user-456"}, nil)

		req := httptest.NewRequest(http.MethodDelete, "/projects/project-123", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "only project creator or leader can delete project", response["error"])
	})

}

func TestUpdateProject(t *testing.T) {
	t.Run("Update Successful", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)

		db, sqlMock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		handler := &routes.ProjectHandler{Repo: mockRepo, DB: db}

		app.Use(withUserID("user-123"))
		app.Put("/projects/:id", handler.UpdateProject)

		mockRepo.On("GetOwnerByProjectID", mock.Anything, "project-123").Return("user-123", nil)
		mockRepo.On("GetProjectByID", mock.Anything, "project-123").
			Return(&models.CreateProject{ID: "project-123", Name: "Old Name"}, nil)

		sqlMock.ExpectExec("UPDATE .*").WillReturnResult(sqlmock.NewResult(1, 1))

		body := `{"name":"Updated Name"}`
		req := httptest.NewRequest(http.MethodPut, "/projects/project-123", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Project updated successfully", response["message"])

		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Project Id is not provided", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Patch("/projects/:id", handler.UpdateProject)

		req := httptest.NewRequest(http.MethodPatch, "/projects/undefined", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "missing project ID", response["error"])
	})

	t.Run("Invalid body request", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Patch("/projects/:id", handler.UpdateProject)

		req := httptest.NewRequest(http.MethodPatch, "/projects/project-123", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "invalid request body", response["error"])
	})

	t.Run("Get Project Owner Failed", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Patch("/projects/:id", handler.UpdateProject)

		mockRepo.On("GetOwnerByProjectID", mock.Anything, "project-123").Return("", errors.New("failed to get project owner"))

		req := httptest.NewRequest(http.MethodPatch, "/projects/project-123", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "project not found", response["error"])
	})

	t.Run("Get Leader Failed", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Patch("/projects/:id", handler.UpdateProject)

		mockRepo.On("GetOwnerByProjectID", mock.Anything, "project-123").
			Return("user-456", nil)

		mockRepo.On("GetLeaderByProjectID", mock.Anything, "project-123").
			Return("", errors.New("leader fetch failed"))

		req := httptest.NewRequest(http.MethodPatch, "/projects/project-123", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "project not found", response["error"])
	})

	t.Run("User is neither Owner nor Leader", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Patch("/projects/:id", handler.UpdateProject)

		mockRepo.On("GetOwnerByProjectID", mock.Anything, "project-123").
			Return("user-456", nil)

		mockRepo.On("GetLeaderByProjectID", mock.Anything, "project-123").
			Return("user-789", nil)

		req := httptest.NewRequest(http.MethodPatch, "/projects/project-123", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "you are not authorized to update this project", response["error"])
	})

	t.Run("GetProjectByID failed", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		handler := &routes.ProjectHandler{Repo: mockRepo}

		app.Use(withUserID("user-123"))
		app.Patch("/projects/:id", handler.UpdateProject)

		mockRepo.On("GetOwnerByProjectID", mock.Anything, "project-123").
			Return("user-123", nil)

		mockRepo.On("GetProjectByID", mock.Anything, "project-123").
			Return(db.Project{}, errors.New("db error"))

		req := httptest.NewRequest(http.MethodPatch, "/projects/project-123", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to fetch project", response["error"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("Add/Remove members loop executed correctly with struct body", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockDb, sqlMock, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDb.Close()

		handler := &routes.ProjectHandler{Repo: mockRepo, DB: mockDb}
		app.Use(withUserID("user-123"))
		app.Patch("/projects/:id", handler.UpdateProject)

		mockRepo.On("GetOwnerByProjectID", mock.Anything, "project-123").
			Return("user-123", nil)
		mockRepo.On("GetProjectByID", mock.Anything, "project-123").
			Return(db.Project{ID: "project-123"}, nil)

		mockRepo.On("AddMemberToProject", mock.Anything, "project-123", "userA").Return(nil).Once()
		mockRepo.On("AddMemberToProject", mock.Anything, "project-123", "userB").Return(nil).Once()
		mockRepo.On("RemoveMemberFromProject", mock.Anything, "project-123", "userX").Return(nil).Once()

		sqlMock.ExpectExec("UPDATE .*").WillReturnResult(sqlmock.NewResult(1, 1))

		addMembers := []string{"userA", "userB"}
		removeMembers := []string{"userX"}
		name := "New Name"
		bodyStruct := models.EditProject{
			AddMember:    &addMembers,
			RemoveMember: &removeMembers,
			Name:         &name,
		}

		jsonBody, err := json.Marshal(bodyStruct)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/projects/project-123", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Project updated successfully", response["message"])

		mockRepo.AssertExpectations(t)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("DB ExecContext returns error", func(t *testing.T) {
		app := fiber.New()

		mockRepo := new(mocks.MockProjectRepo)
		mockDb, sqlMock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDb.Close()

		handler := &routes.ProjectHandler{Repo: mockRepo, DB: mockDb}
		app.Use(withUserID("user-123"))
		app.Patch("/projects/:id", handler.UpdateProject)

		mockRepo.On("GetOwnerByProjectID", mock.Anything, "project-123").Return("user-123", nil)
		mockRepo.On("GetProjectByID", mock.Anything, "project-123").Return(db.Project{ID: "project-123"}, nil)

		emptySlice := []string{}
		bodyStruct := models.EditProject{
			AddMember:    &emptySlice,
			RemoveMember: &emptySlice,
			Name:         nil,
		}
		jsonBody, err := json.Marshal(bodyStruct)
		require.NoError(t, err)

		sqlMock.ExpectExec("UPDATE .*").WillReturnError(fmt.Errorf("some exec error"))

		req := httptest.NewRequest(http.MethodPatch, "/projects/project-123", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "failed to update project", response["error"])

		mockRepo.AssertExpectations(t)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

}
