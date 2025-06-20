package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/auth"
)

func SetUpAuthRoutes(api fiber.Router, h *auth.AuthHandler) {
	api.Post("/login", h.Login)
	api.Post("/register", h.Register)
}
