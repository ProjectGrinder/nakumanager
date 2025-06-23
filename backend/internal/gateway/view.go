package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpViewRoutes(api fiber.Router, h *routes.ViewHandler) {
	api.Post("/views", h.CreateView)
	api.Get("/views/:id", h.GetViewsByUserID)
	api.Delete("/views/:id", h.DeleteView)
}
