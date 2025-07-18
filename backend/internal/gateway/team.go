package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpTeamRoutes(api fiber.Router, h *routes.TeamHandler) {
	api.Post("/teams", h.CreateTeam)
	api.Get("/teams", h.GetTeamsByUserID)
	api.Patch("/teams/:id", h.UpdateTeam)
	api.Delete("/teams/:id", h.DeleteTeam)

}
