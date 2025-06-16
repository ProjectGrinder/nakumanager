package ws

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func UpdateViewHandler(c *websocket.Conn) {
	log.Println("update view")

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			log.Println("WebSocket disconnected:", err)
			break
		}
	}
}
