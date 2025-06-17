package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

func UpdateProjectHandler(c *websocket.Conn, data json.RawMessage) {
	c.WriteMessage(websocket.TextMessage, []byte("project updated"))

}
