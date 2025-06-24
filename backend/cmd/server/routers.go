package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/auth"
	"github.com/nack098/nakumanager/internal/db"
	"github.com/nack098/nakumanager/internal/gateway"
	"github.com/nack098/nakumanager/internal/repositories"
	"github.com/nack098/nakumanager/internal/routes"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	// ก่อนที่จะเรียก handler ตัวจริง
	method := c.Method()
	path := c.Path()
	ip := c.IP()

	fmt.Printf("[LOG] %s %s - from IP: %s\n", method, path, ip)

	// เรียก handler ตัวจริงต่อ
	return c.Next()
}

func SetUpRouters(app *fiber.App, conn *sql.DB) {
	queries := db.New(conn)
	userRepo := repositories.NewUserRepository(queries)
	workspaceRepo := repositories.NewWorkspaceRepository(queries)
	teamRepo := repositories.NewTeamRepository(queries)
	projectRepo := repositories.NewProjectRepository(queries)
	issueRepo := repositories.NewIssueRepository(queries)
	viewRepo := repositories.NewViewRepository(conn)

	authHandler := auth.NewAuthHandler(userRepo)
	workspaceHandler := routes.NewWorkspaceHandler(workspaceRepo, userRepo)
	teamHandler := routes.NewTeamHandler(teamRepo, workspaceRepo)
	projectHandler := routes.NewProjectHandler(projectRepo, teamRepo)
	issueHandler := routes.NewIssueHandler(issueRepo, teamRepo, projectRepo)
	viewHandler := routes.NewViewHandler(viewRepo)

	api := app.Group("/api")

	api.Use(LoggerMiddleware)

	gateway.SetUpAuthRoutes(api, authHandler)

	private := api.Group("/")
	private.Use(authHandler.AuthRequired)

	gateway.SetUpWorkspaceRoutes(private, workspaceHandler)
	gateway.SetUpTeamRoutes(private, teamHandler)
	gateway.SetUpProjectsRoutes(private, projectHandler)
	gateway.SetUpIssueRoutes(private, issueHandler)
	gateway.SetUpViewRoutes(private, viewHandler)

	// wsGroup := app.Group("/ws", ws.WebSocketMiddleware)
	// wsGroup.Get("/", websocket.New(ws.CentralWebSocketHandler))

}
