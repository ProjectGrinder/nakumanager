package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetUpViewRoutes(api fiber.Router){
	api.Post("/views", CreateView)
	api.Get("/views/:id", GetViewsByUserID)
	api.Delete("/views/:id", DeleteView)
}


func CreateView(c *fiber.Ctx) error {
	return c.SendString("Hello From Create View!")
}

func GetViewsByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Views!")
}

func DeleteView(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete View!")
}