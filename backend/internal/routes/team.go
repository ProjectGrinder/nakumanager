package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/repositories"
)

type TeamHandler struct {
	Repo          repositories.TeamRepository
	WorkspaceRepo repositories.WorkspaceRepository
}

func NewTeamHandler(repo repositories.TeamRepository, workspaceRepo repositories.WorkspaceRepository) *TeamHandler {
	return &TeamHandler{
		Repo:          repo,
		WorkspaceRepo: workspaceRepo,
	}
}

func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	return c.SendString("Hello From Create Team!")
}

func (h *TeamHandler) GetTeamsByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Teams!")
}

func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete Team!")
}
