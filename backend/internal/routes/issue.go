package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/repositories"
)

type IssueHandler struct {
	Repo repositories.IssueRepository
}

func NewIssueHandler(repo repositories.IssueRepository) *IssueHandler {
	return &IssueHandler{Repo: repo}
}


//TODO : implement All the functions
func (h *IssueHandler)CreateIssue(c *fiber.Ctx) error {
	uuid := uuid.NewString()
	return c.Status(200).JSON(uuid)
}

func (h *IssueHandler)GetIssuesByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	return c.Status(200).JSON(userID)
}

func (h *IssueHandler)DeleteIssue(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	return c.Status(200).JSON(userID)
}
