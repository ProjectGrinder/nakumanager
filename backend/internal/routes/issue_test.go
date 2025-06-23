package routes_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/routes"
	mocks "github.com/nack098/nakumanager/internal/routes/test/mock_repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateIssue_Success(t *testing.T) {
	app := fiber.New()

	mockTeamRepo := new(mocks.MockTeamRepo)
	mockProjectRepo := new(mocks.MockProjectRepo)
	mockIssueRepo := new(mocks.MockIssueRepo)

	handler := &routes.IssueHandler{
		TeamRepo:    mockTeamRepo,
		ProjectRepo: mockProjectRepo,
		Repo:        mockIssueRepo,
	}

	app.Post("/api/issues", func(c *fiber.Ctx) error {
		c.Locals("userID", "user-123")
		return handler.CreateIssue(c)
	})

	payload := models.IssueCreate{
		Title:   "Test Issue",
		TeamID:  "team-001",
		Status:  "todo",
		OwnerID: "",
	}
	body, _ := json.Marshal(payload)

	mockTeamRepo.On("IsTeamExists", mock.Anything, "team-001").Return(true, nil)
	mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-001", "user-123").Return(true, nil)
	mockIssueRepo.On("CreateIssue", mock.Anything, mock.AnythingOfType("db.CreateIssueParams")).Return(nil)

	req := httptest.NewRequest("POST", "/api/issues", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
	mockIssueRepo.AssertExpectations(t)
}

func TestAddAssigneeToIssue_Success(t *testing.T) {
	app := fiber.New()

	mockIssueRepo := new(mocks.MockIssueRepo)
	handler := routes.NewIssueHandler(mockIssueRepo, nil, nil)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "auth-user-1")
		return c.Next()
	})

	app.Post("/add-assignee", handler.AddAssigneeToIssue)

	mockIssueRepo.On("AddAssigneeToIssue", mock.Anything, db.AddAssigneeToIssueParams{
		IssueID: "issue-123",
		UserID:  "user-456",
	}).Return(nil)

	reqBody := `{"issue_id":"issue-123","user_id":"user-456"}`

	req := httptest.NewRequest("POST", "/add-assignee", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Logf("Error during app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Logf("Expected status %d but got %d", fiber.StatusOK, resp.StatusCode)
	}

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockIssueRepo.AssertExpectations(t)
}

func TestRemoveAssigneeFromIssue_Success(t *testing.T) {
	app := fiber.New()

	mockIssueRepo := new(mocks.MockIssueRepo)
	handler := routes.NewIssueHandler(mockIssueRepo, nil, nil)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "auth-user-1")
		return c.Next()
	})

	app.Post("/remove-assignee", handler.RemoveAssigneeFromIssue)

	mockIssueRepo.On("RemoveAssigneeFromIssue", mock.Anything, db.RemoveAssigneeFromIssueParams{
		IssueID: "issue-123",
		UserID:  "user-456",
	}).Return(nil)

	reqBody := `{"issue_id":"issue-123","user_id":"user-456"}`

	req := httptest.NewRequest("POST", "/remove-assignee", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockIssueRepo.AssertExpectations(t)
}

func TestGetIssuesByUserID_Success(t *testing.T) {
	app := fiber.New()

	mockIssueRepo := new(mocks.MockIssueRepo)
	handler := routes.NewIssueHandler(mockIssueRepo, nil, nil)

	app.Get("/issues/:id", handler.GetIssuesByUserID)

	expectedIssues := []db.Issue{
		{ID: "issue-1", Title: "Issue One", Status: "todo"},
		{ID: "issue-2", Title: "Issue Two", Status: "done"},
	}

	mockIssueRepo.On("GetIssueByUserID", mock.Anything, "user-123").
		Return(expectedIssues, nil).
		Once()

	req := httptest.NewRequest("GET", "/issues/user-123", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var actualIssues []db.Issue
	err = json.Unmarshal(body, &actualIssues)
	assert.NoError(t, err)
	assert.Equal(t, expectedIssues, actualIssues)

	mockIssueRepo.AssertExpectations(t)
}

func TestDeleteIssue_Success(t *testing.T) {
	app := fiber.New()

	mockIssueRepo := new(mocks.MockIssueRepo)
	handler := routes.NewIssueHandler(mockIssueRepo, nil, nil)

	userID := "user-456"
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})

	app.Delete("/issues/:id", handler.DeleteIssue)

	issueID := "issue-123"

	mockIssueRepo.On("GetIssueByID", mock.Anything, issueID).
		Return(db.Issue{ID: issueID, OwnerID: userID}, nil).
		Once()

	mockIssueRepo.On("DeleteIssue", mock.Anything, issueID).
		Return(nil).
		Once()

	req := httptest.NewRequest("DELETE", "/issues/"+issueID, nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockIssueRepo.AssertExpectations(t)
}
