package routes

import (
	"github.com/gofiber/fiber/v2"
)

func CreateView(c *fiber.Ctx) error {
	return c.SendString("Hello From Create View!")
}

func GetViewsByID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Views!")
}

func DeleteView(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete View!")
}