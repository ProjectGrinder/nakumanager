package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpIssueRoutes(api fiber.Router, h *routes.IssueHandler) {
	api.Post("/issues", h.CreateIssue)
	api.Post("/issues/assignee", h.AddAssigneeToIssue)
	api.Delete("/issues/assignee", h.RemoveAssigneeFromIssue)
	api.Get("/issues/:id", h.GetIssuesByUserID)
	api.Delete("/issues/:id", h.DeleteIssue)

}
