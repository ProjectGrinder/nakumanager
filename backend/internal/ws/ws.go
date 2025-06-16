package ws

import (
	"log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WebSocketMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WebSocketHandler(c *websocket.Conn) {
	userID := c.Params("id")
	log.Println("Connected user:", userID)

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err := c.WriteMessage(mt, msg); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
