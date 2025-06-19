package main

import (
	
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	SetUpRouters(app)

	app.Listen(":3000")
}
