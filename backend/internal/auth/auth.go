package auth

import(
	"github.com/gofiber/fiber/v2"
)


func Login(c *fiber.Ctx) error {
	return c.SendString("Hello From Login!")
}

func Register(c *fiber.Ctx) error {
	return c.SendString("Hello From Register!")
}


