package routes_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/db"
	"github.com/nack098/nakumanager/internal/routes"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateView_SingleGroupBy(t *testing.T) {
	app := fiber.New()

	mockRepo := new(mocks.MockViewRepo)
	handler := &routes.ViewHandler{Repo: mockRepo}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "user123")
		return c.Next()
	})

	app.Post("/views", handler.CreateView)

	groupBys := []string{"status"}

	body := map[string]interface{}{
		"name":      "Single Group View",
		"team_id":   "team123",
		"group_bys": groupBys,
	}
	bodyBytes, _ := json.Marshal(body)

	
	mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

	mockRepo.On("GetGroupedIssues", mock.Anything, "team123", groupBys).
		Return([]map[string]interface{}{
			{"status": sql.NullString{Valid: true, String: "todo"}},
		}, nil)

	mockRepo.On("ListIssuesByGroupFilters", mock.Anything, "team123", map[string]string{
		"status": "todo",
	}).Return([]db.Issue{{ID: "issue1"}}, nil)

	mockRepo.On("AddIssueToView", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestCreateView_MultiGroupBy(t *testing.T) {
	app := fiber.New()

	mockRepo := new(mocks.MockViewRepo)
	handler := &routes.ViewHandler{Repo: mockRepo}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "user123")
		return c.Next()
	})

	app.Post("/views", handler.CreateView)

	groupBys := []string{"status", "priority"}

	body := map[string]interface{}{
		"name":      "Multi Group View",
		"team_id":   "team123",
		"group_bys": groupBys,
	}
	bodyBytes, _ := json.Marshal(body)

	mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
	for _, g := range groupBys {
		mockRepo.On("AddGroupByToView", mock.Anything, mock.MatchedBy(func(p db.AddGroupByToViewParams) bool {
			return p.GroupBy == g
		})).Return(nil)
	}

	mockRepo.On("GetGroupedIssues", mock.Anything, "team123", groupBys).
		Return([]map[string]interface{}{
			{"status": sql.NullString{Valid: true, String: "todo"}, "priority": sql.NullString{Valid: true, String: "High"}},
			{"status": sql.NullString{Valid: true, String: "Done"}, "priority": sql.NullString{Valid: true, String: "Low"}},
		}, nil)

	mockRepo.On("ListIssuesByGroupFilters", mock.Anything, "team123", map[string]string{
		"status": "todo", "priority": "High",
	}).Return([]db.Issue{{ID: "issue1"}}, nil)

	mockRepo.On("ListIssuesByGroupFilters", mock.Anything, "team123", map[string]string{
		"status": "Done", "priority": "Low",
	}).Return([]db.Issue{{ID: "issue2"}}, nil)

	mockRepo.On("AddIssueToView", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestGetViewsByGroupBy_Success(t *testing.T) {
	app := fiber.New()

	mockRepo := new(mocks.MockViewRepo)
	handler := &routes.ViewHandler{Repo: mockRepo}

	app.Post("/views/groupby", handler.GetViewsByGroupBy)

	viewID := "view123"

	body := map[string]interface{}{
		"view_id": viewID,
	}
	bodyBytes, _ := json.Marshal(body)
	mockRepo.On("ListGroupByViewID", mock.Anything, viewID).
		Run(func(args mock.Arguments) {
			log.Printf("[MOCK] ListGroupByViewID called with viewID=%v", args.Get(1))
		}).
		Return([]string{"status", "priority"}, nil)

	mockRepo.On("ListIssuesByViewID", mock.Anything, viewID).
		Run(func(args mock.Arguments) {
			log.Printf("[MOCK] ListIssuesByViewID called with viewID=%v", args.Get(1))
		}).
		Return([]db.Issue{
			{
				ID:       "issue1",
				Status:   "todo",
				Priority: sql.NullString{Valid: true, String: "High"},
			},
			{
				ID:       "issue2",
				Status:   "Done",
				Priority: sql.NullString{Valid: true, String: "Low"},
			},
		}, nil)

	req := httptest.NewRequest(http.MethodPost, "/views/groupby", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	t.Logf("RESPONSE BODY: %s", string(respBody))


	var groups []map[string]interface{}
	err = json.Unmarshal(respBody, &groups)
	assert.NoError(t, err)


	assert.Len(t, groups, 2)

	groupKeys := map[string]int{}
	for _, group := range groups {
		status := group["status"]
		priority := group["priority"]
		key := fmt.Sprintf("%v|%v", status, priority)
		groupKeys[key]++
	}

	assert.Equal(t, 1, groupKeys["todo|High"])
	assert.Equal(t, 1, groupKeys["Done|Low"])
}
