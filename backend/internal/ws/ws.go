package ws

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/auth"
)

func WebSocketMiddleware(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	token := c.Cookies("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authentication token",
		})
	}

	//TODO: fix return token
	if _ ,err := auth.VerifyToken(token); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	c.Locals("token", token)
	return c.Next()
}

func CentralWebSocketHandler(c *websocket.Conn) {
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var message struct {
			Event string          `json:"event"`
			Data  json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("json error:", err)
			continue
		}

		switch message.Event {
		case "update_workspace":
			UpdateWorkspaceHandler(c, message.Data)
		case "update_project":
			UpdateProjectHandler(c, message.Data)
		case "update_issue":
			UpdateIssueHandler(c, message.Data)
		case "update_view":
			UpdateViewHandler(c, message.Data)
		case "update_team":
			UpdateTeamHandler(c, message.Data)
		default:
			log.Println("unknown event:", message.Event)
		}

	}
}
