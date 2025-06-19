package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetUpWorkspaceRoutes(api fiber.Router) {
	api.Get("/workspace/:id", GetWorkspacesByUserID)
	api.Post("/workspace", CreateWorkspace)
	api.Delete("/workspace/:id", DeleteWorkspace)
}

func CreateWorkspace(c *fiber.Ctx) error {
	return c.SendString("Hello From Create Workspace!")
}

func GetWorkspacesByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Workspaces!")
}
func DeleteWorkspace(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete Workspace!")
}
