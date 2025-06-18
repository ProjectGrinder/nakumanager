package routes

import (
	"github.com/gofiber/fiber/v2"
)


func SetUpIssueRoutes(api fiber.Router) {
	api.Post("/issues", CreateIssue)
	api.Get("/issues/:id", GetIssuesByUserID)
	api.Delete("/issues/:id", DeleteIssue)
}

func CreateIssue(c *fiber.Ctx) error {
	return c.SendString("Hello From Create Issue!")
}

func GetIssuesByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Issues!")
}

func DeleteIssue(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete Issue!")
}