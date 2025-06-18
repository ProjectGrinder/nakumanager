package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/auth"
	"github.com/nack098/nakumanager/internal/routes"
	"github.com/nack098/nakumanager/internal/ws"
)

func main() {
	app := fiber.New()

	api := app.Group("/api")

	// Auth
	api.Post("/login", auth.Login)
	api.Post("/register", auth.Register)

	// User
	api.Post("/users", routes.CreateUser)
	api.Get("/users", routes.GetAllUsers)
	api.Get("/user", routes.GetUserByID)
	api.Delete("/users", routes.DeleteUser)

	//workspace
	api.Post("/workspaces", routes.CreateWorkspace)
	api.Get("/workspaces", routes.GetWorkspacesByUserID)
	api.Delete("/workspaces", routes.DeleteWorkspace)

	// Project
	api.Post("/projects", routes.CreateProject)
	api.Get("/projects", routes.GetProjectsByUserID)
	api.Delete("/projects", routes.DeleteProject)

	// Team
	api.Post("/teams", routes.CreateTeam)
	api.Get("/teams", routes.GetTeamsByUserID)
	api.Delete("/teams", routes.DeleteTeam)

	// Task
	api.Post("/issues", routes.CreateIssue)
	api.Get("/issues", routes.GetIssuesByUserID)
	api.Delete("/issues", routes.DeleteIssue)

	// View
	api.Post("/views", routes.CreateView)
	api.Get("/views", routes.GetViewsByUserID)
	api.Delete("/views", routes.DeleteView)

	// Websocket
	wsGroup := app.Group("/ws", ws.WebSocketMiddleware)
	wsGroup.Get("/", websocket.New(ws.CentralWebSocketHandler))

	app.Listen(":3000")
}
