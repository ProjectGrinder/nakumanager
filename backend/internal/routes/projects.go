package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/repositories"
)

type ProjectHandler struct {
	Repo repositories.ProjectRepository
}

func NewProjectHandler(repo repositories.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{Repo: repo}
}


func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	return c.SendString("Hello From Create Project!")
}

func (h *ProjectHandler) GetProjectsByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Projects!")
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete Project!")
}
