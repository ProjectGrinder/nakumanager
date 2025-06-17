package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

func UpdateViewHandler(c *websocket.Conn, data json.RawMessage) {
	c.WriteMessage(websocket.TextMessage, []byte("view updated"))

}
