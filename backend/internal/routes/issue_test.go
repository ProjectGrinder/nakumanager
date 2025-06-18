package routes_test

import (
	"net/http/httptest"
	"testing"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)

func setUpIssueApp() *fiber.App {
	app := fiber.New()
	app.Get("/issues", routes.GetIssuesByUserID)
	app.Post("/issues", routes.CreateIssue)
	app.Delete("/issues", routes.DeleteIssue)
	return app
}

func TestCreateIssue(t *testing.T) {
	app := setUpIssueApp()
	req := httptest.NewRequest("POST", "/issues", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Create Issue!", string(body))
}

func TestGetIssueByUserID(t *testing.T) {
	app := setUpIssueApp()
	req := httptest.NewRequest("GET", "/issues", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Get Issues!", string(body))
}

func TestDeleteIssue(t *testing.T) {
	app := setUpIssueApp()
	req := httptest.NewRequest("DELETE", "/issues", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Delete Issue!", string(body))
}
