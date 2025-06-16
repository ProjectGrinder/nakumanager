package ws

import(
	"log"
	"github.com/gofiber/contrib/websocket"
)


func UpdateTeamHandler(c *websocket.Conn) {
	log.Println("update team")

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			log.Println("WebSocket disconnected:", err)
			break
		}
	}
}