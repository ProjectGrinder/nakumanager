package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpWorkspaceRoutes(api fiber.Router, h *routes.WorkspaceHandler) {
	api.Get("/workspace", h.GetWorkspacesByUserID)
	api.Post("/workspace", h.CreateWorkspace)
	api.Post("/workspace/:workspaceid/members", h.AddMemberToWorkspace)
	api.Delete("/workspace/:workspaceid/members", h.RemoveMemberFromWorkspace)
	api.Post("/workspace/:workspaceid/rename", h.RenameWorkSpace)
	api.Delete("/workspace/:workspaceid", h.DeleteWorkspace)

}
