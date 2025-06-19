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
	routes.SetUpUserRoutes(api)
	routes.SetUpWorkspaceRoutes(api)
	routes.SetUpProjectsRoutes(api)
	routes.SetUpTeamRoutes(api)
	routes.SetUpIssueRoutes(api)
	routes.SetUpViewRoutes(api)

	wsGroup := app.Group("/ws", ws.WebSocketMiddleware)
	wsGroup.Get("/", websocket.New(ws.CentralWebSocketHandler))
}
