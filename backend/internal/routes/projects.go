package routes

import(
	"github.com/gofiber/fiber/v2"
)

func CreateProject(c *fiber.Ctx) error {
	return c.SendString("Hello From Create Project!")
}

func GetProjectsByID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Projects!")
}

func DeleteProject(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete Project!")
}


