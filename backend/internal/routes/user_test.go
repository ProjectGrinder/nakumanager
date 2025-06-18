package routes_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/stretchr/testify/assert"
)

func setUpUserApp() *fiber.App {
	app := fiber.New()
	app.Post("/users", routes.CreateUser)
	app.Get("/users", routes.GetAllUsers)
	app.Get("/user", routes.GetUserByID) // รับ query param: /user?id=123
	app.Delete("/users", routes.DeleteUser)
	return app
}

func TestCreateUser(t *testing.T) {
	app := setUpUserApp()
	req := httptest.NewRequest("POST", "/users", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Create User!", string(body))
}

func TestGetAllUsers(t *testing.T) {
	app := setUpUserApp()
	req := httptest.NewRequest("GET", "/users", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Get Users!", string(body))
}

func TestGetUserByID(t *testing.T) {
	app := setUpUserApp()
	req := httptest.NewRequest("GET", "/user?id=123", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Get User! ID: 123", string(body))
}

func TestDeleteUser(t *testing.T) {
	app := setUpUserApp()
	req := httptest.NewRequest("DELETE", "/users", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello From Delete User!", string(body))
}
