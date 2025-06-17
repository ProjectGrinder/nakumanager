package routes

import (
	"github.com/gofiber/fiber/v2"
)


func CreateIssue(c *fiber.Ctx) error {
	return c.SendString("Hello From Create Issue!")
}

func GetIssuesByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Issues!")
}

func DeleteIssue(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete Issue!")
}