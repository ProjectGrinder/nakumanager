package ws

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func UpdateWorkspaceHandler(c *websocket.Conn) {
	log.Println("update workspace")

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			log.Println("WebSocket disconnected:", err)
			break
		}
	}
}
