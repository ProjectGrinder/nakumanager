package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetUpTeamRoutes(api fiber.Router) {
	api.Post("/teams", CreateTeam)
	api.Get("/teams/:id", GetTeamsByUserID)
	api.Delete("/teams/:id", DeleteTeam)
}

func CreateTeam(c *fiber.Ctx) error {
	return c.SendString("Hello From Create Team!")
}

func GetTeamsByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Teams!")
}

func DeleteTeam(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete Team!")
}