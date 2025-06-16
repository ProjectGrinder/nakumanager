package ws

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func UpdateIssueHandler(c *websocket.Conn) {
	log.Println("update issue")
	
	for {
		if _, _, err := c.ReadMessage(); err != nil {
			log.Println("WebSocket disconnected:", err)
			break
		}
	}
}
