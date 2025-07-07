package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpTeamRoutes(api fiber.Router, h *routes.TeamHandler) {
	api.Post("/teams", h.CreateTeam)
	api.Get("/teams", h.GetTeamsByUserID)
	api.Post("/teams/:id/members", h.AddMemberToTeam)
	api.Delete("/teams/:id/members", h.RemoveMemberFromTeam)
	api.Post("/teams/:id/rename", h.RenameTeam)
	api.Post("/teams/:id/leader", h.SetTeamLeader)
	api.Delete("/teams/:id", h.DeleteTeam)

}
