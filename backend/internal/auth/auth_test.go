package auth_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {

	app := fiber.New()
	app.Post("/login", auth.Login)
	req := httptest.NewRequest("POST", "/login", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestRegister(t *testing.T) {

	app := fiber.New()
	app.Post("/register", auth.Register)
	req := httptest.NewRequest("POST", "/register", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

}
