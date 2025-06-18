package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/nack098/nakumanager/internal/auth"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/nack098/nakumanager/internal/ws"
)

func SetUpRouters(app *fiber.App) {
	api := app.Group("/api")

	auth.SetUpAuthRoutes(api)

	private := api.Group("/")
	private.Use(auth.AuthRequired)

	routes.SetUpUserRoutes(private)
	routes.SetUpWorkspaceRoutes(private)
	routes.SetUpProjectsRoutes(private)
	routes.SetUpTeamRoutes(private)
	routes.SetUpIssueRoutes(private)
	routes.SetUpViewRoutes(private)

	wsGroup := app.Group("/ws", ws.WebSocketMiddleware)
	wsGroup.Get("/", websocket.New(ws.CentralWebSocketHandler))
}

