package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpUserRoutes(api fiber.Router, h *routes.UserHandler) {
	api.Post("/issues", h.CreateUser)
	api.Get("/issues", h.GetAllUsers)
	api.Get("/issues/:id", h.GetUserByID)
	api.Delete("/issues/:id", h.DeleteUser)
}
