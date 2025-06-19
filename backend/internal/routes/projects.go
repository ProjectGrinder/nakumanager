package routes

import(
	"github.com/gofiber/fiber/v2"
)

func SetUpProjectsRoutes(api fiber.Router) {
	api.Get("/projects/:id", GetProjectsByUserID)
	api.Post("/projects", CreateProject)
	api.Delete("/projects/:id", DeleteProject)
}

func CreateProject(c *fiber.Ctx) error {
	return c.SendString("Hello From Create Project!")
}

func GetProjectsByUserID(c *fiber.Ctx) error {
	return c.SendString("Hello From Get Projects!")
}

func DeleteProject(c *fiber.Ctx) error {
	return c.SendString("Hello From Delete Project!")
}


