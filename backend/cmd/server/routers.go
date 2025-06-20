package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/auth"
	"github.com/nack098/nakumanager/internal/db"
	"github.com/nack098/nakumanager/internal/gateway"
	"github.com/nack098/nakumanager/internal/repositories"
	"github.com/nack098/nakumanager/internal/routes"
)

func SetUpRouters(app *fiber.App, conn *sql.DB) {
	queries := db.New(conn)
	userRepo := repositories.NewUserRepository(queries)
	workspaceRepo := repositories.NewWorkspaceRepository(queries)

	authHandler := auth.NewAuthHandler(userRepo)
	workspaceHandler := routes.NewWorkspaceHandler(workspaceRepo, userRepo)

	api := app.Group("/api")

	gateway.SetUpAuthRoutes(api, authHandler)

	private := api.Group("/")
	private.Use(authHandler.AuthRequired)

	gateway.SetUpWorkspaceRoutes(private, workspaceHandler)

}

// routes.SetUpUserRoutes(private)
// routes.SetUpProjectsRoutes(private)
// routes.SetUpTeamRoutes(private)
// routes.SetUpIssueRoutes(private)
// routes.SetUpViewRoutes(private)

// wsGroup := app.Group("/ws", ws.WebSocketMiddleware)
// wsGroup.Get("/", websocket.New(ws.CentralWebSocketHandler))
