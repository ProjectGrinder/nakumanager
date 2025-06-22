package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpProjectsRoutes(api fiber.Router, h *routes.ProjectHandler) {
	api.Get("/projects", h.GetProjectsByUserID)
	api.Post("/projects", h.CreateProject)
	api.Post("/projects/members", h.AddMemberToProject)
	api.Delete("/projects/members", h.RemoveProjectMembers)
	api.Delete("/projects/:id", h.DeleteProject)
}
