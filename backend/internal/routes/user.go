package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetUpUserRoutes(api fiber.Router) {
	api.Post("/users", CreateUser)
	api.Get("/users", GetAllUsers)
	api.Get("/users/:id", GetUserByID)
	api.Delete("/users/:id", DeleteUser)
}

func CreateUser(c *fiber.Ctx) error {
	return c.SendString("Hello From Create User!")
}

func GetAllUsers(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Users!")
}

func GetUserByID(c *fiber.Ctx) error {
	uId := c.Query("id")
	return c.SendString("Hello From Get User! ID: " + uId)
}


func DeleteUser(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete User!")
}