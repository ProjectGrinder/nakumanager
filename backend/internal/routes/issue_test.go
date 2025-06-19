package routes_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/auth"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	api.Use("/issues", auth.AuthRequired)
	routes.SetUpIssueRoutes(api)

	return app
}

func TestCreateIssue_Success(t *testing.T) {
	app := fiber.New()

	app.Post("/issues", func(c *fiber.Ctx) error {
		c.Locals("userID", "mockedID")
		return routes.CreateIssue(c)
	})

	payload := `{"title":"Test Issue","teamId":1112}`

	req := httptest.NewRequest("POST", "/issues", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

}

func TestCreateIssue_InvalidPayload(t *testing.T) {
	app := fiber.New()

	app.Post("/issues", func(c *fiber.Ctx) error {
		c.Locals("userID", "mockedID")
		return routes.CreateIssue(c)
	})

	payload := `{"title":"Test Issue"}`
	req := httptest.NewRequest("POST", "/issues", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCreateIssue_Unauthorized(t *testing.T) {
	app := fiber.New()

	app.Post("/issues", func(c *fiber.Ctx) error {
		return routes.CreateIssue(c)
	})

	payload := `{"title":"Test Issue","teamId":1112}`

	req := httptest.NewRequest("POST", "/issues", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
func TestCreateIssue_InvalidBody(t *testing.T) {
	app := fiber.New()
	app.Post("/api/issues", routes.CreateIssue)

	req := httptest.NewRequest("POST", "/api/issues", strings.NewReader("not-a-json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "invalid request body")
}

func TestGetIssuesByUserID_Success(t *testing.T) {
	app := fiber.New()

	app.Get("/api/issues/:id", routes.GetIssuesByUserID)

	issue1 := models.Issue{
		ID:       "1",
		Title:    "Issue 1",
		Assignee: []string{"user1", "user2"},
	}
	issue2 := models.Issue{
		ID:       "2",
		Title:    "Issue 2",
		Assignee: []string{"user3"},
	}

	routes.Issues = map[string]models.Issue{
		"1": issue1,
		"2": issue2,
	}

	req := httptest.NewRequest("GET", "/api/issues/user1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "Issue 1")
}

func TestGetIssuesByUserID_Unsuccess(t *testing.T) {
	app := fiber.New()

	app.Get("/api/issues/:id", routes.GetIssuesByUserID)

	issue1 := models.Issue{
		ID:       "1",
		Title:    "Issue 1",
		Assignee: []string{"user1", "user2"},
	}
	issue2 := models.Issue{
		ID:       "2",
		Title:    "Issue 2",
		Assignee: []string{"user3"},
	}

	routes.Issues = map[string]models.Issue{
		"1": issue1,
		"2": issue2,
	}

	req := httptest.NewRequest("GET", "/api/issues/user4", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "null", strings.TrimSpace(string(body)))
}

func TestDeleteIssue_Success(t *testing.T) {
	app := fiber.New()
	app.Delete("/issues/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "user-123")
		return routes.DeleteIssue(c)
	})

	issueID := "issue-1"
	routes.Issues = map[string]models.Issue{
		issueID: {
			ID:       issueID,
			OwnerID:  "user-123",
			Assignee: []string{"user-456"},
		},
	}

	req := httptest.NewRequest("DELETE", "/issues/"+issueID, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Issue deleted successfully")
}

func TestDeleteIssue_IssueNotExists(t *testing.T) {
	app := fiber.New()
	app.Delete("/issues/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "user-123")
		return routes.DeleteIssue(c)
	})

	issueID := "issue-1"
	routes.Issues = map[string]models.Issue{
		issueID: {
			ID:       issueID,
			OwnerID:  "user-123",
			Assignee: []string{"user-456"},
		},
	}

	req2 := httptest.NewRequest("DELETE", "/issues/not-exist", nil)
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusNotFound, resp2.StatusCode)
}

func TestDeleteIssue_Forbidden(t *testing.T) {
	issueID := "issue-1"
	routes.Issues[issueID] = models.Issue{
		ID:       issueID,
		OwnerID:  "user-123",
		Assignee: []string{"user-456"},
	}

	userID := "user-789" 

	app := fiber.New()
	app.Delete("/issues/:id", func(c *fiber.Ctx) error {
		if userID != "" {
			c.Locals("userID", userID)
		}
		return routes.DeleteIssue(c)
	})

	req := httptest.NewRequest("DELETE", "/issues/"+issueID, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func TestDeleteIssue_Unauthorized(t *testing.T) {
	issueID := "issue-1"
	routes.Issues[issueID] = models.Issue{
		ID:       issueID,
		OwnerID:  "user-123",
		Assignee: []string{"user-456"},
	}

	app := fiber.New()
	
	app.Delete("/issues/:id", func(c *fiber.Ctx) error {
		
		return routes.DeleteIssue(c)
	})

	req := httptest.NewRequest("DELETE", "/issues/"+issueID, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode) // 401
}
