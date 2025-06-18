package routes_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)

func setUpViewApp() *fiber.App {
	app := fiber.New()
	app.Get("/views", routes.GetViewsByUserID)
	app.Post("/views", routes.CreateView)
	app.Delete("/views", routes.DeleteView)
	return app
}

func TestCreateView(t *testing.T) {
	app := setUpViewApp()
	req := httptest.NewRequest("POST", "/views", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Create View!", string(body))
}


func TestDeleteView(t *testing.T) {
	app := setUpViewApp()
	req := httptest.NewRequest("DELETE", "/views", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Delete View!", string(body))
}

func TestGetViewByUserID(t *testing.T){
	app := setUpViewApp()
	req := httptest.NewRequest("GET", "/views", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Get Views!", string(body))
}