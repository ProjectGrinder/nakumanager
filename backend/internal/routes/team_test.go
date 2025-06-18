package routes_test

import (
	"net/http/httptest"
	"testing"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)


func setUpTeamApp() *fiber.App {
	app := fiber.New()
	app.Post("/teams", routes.CreateTeam)
	app.Get("/teams", routes.GetTeamsByUserID)
	app.Delete("/teams", routes.DeleteTeam)
	return app
}


func TestCreateTeam(t *testing.T) {
	app := setUpTeamApp()
	req := httptest.NewRequest("POST", "/teams", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Create Team!", string(body))
}

func TestGetTeamByUserID(t *testing.T) {
	app := setUpTeamApp()
	req := httptest.NewRequest("GET", "/teams", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Get Teams!", string(body))
}

func TestDeleteTeam(t *testing.T){
	app := setUpTeamApp()
	req := httptest.NewRequest("DELETE", "/teams", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Delete Team!", string(body))
}