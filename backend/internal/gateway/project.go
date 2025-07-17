package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpProjectsRoutes(api fiber.Router, h *routes.ProjectHandler) {
	api.Get("/projects", h.GetProjectsByUserID)
	api.Post("/projects", h.CreateProject)
	api.Patch("/projects/:id", h.UpdateProject)
	api.Delete("/projects/:id", h.DeleteProject)
}
