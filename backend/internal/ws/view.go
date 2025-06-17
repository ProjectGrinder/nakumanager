package ws

import (
	"log"
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

func UpdateViewHandler(c *websocket.Conn, data json.RawMessage) {
	log.Println("update view")

}
