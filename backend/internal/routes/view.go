package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/repositories"
)

type ViewHandler struct {
	Repo repositories.ViewRepository
}

func NewViewHandler(repo repositories.ViewRepository) *ViewHandler {
	return &ViewHandler{Repo: repo}
}

//TODO : implement All the functions
func (h *ViewHandler) CreateView(c *fiber.Ctx) error {
	return c.SendString("Hello From Create View!")
}

func (h *ViewHandler)GetViewsByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Views!")
}

func (h *ViewHandler)DeleteView(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete View!")
}
