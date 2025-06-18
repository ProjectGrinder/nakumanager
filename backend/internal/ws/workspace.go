package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

func UpdateWorkspaceHandler(c *websocket.Conn, data json.RawMessage) {
	c.WriteMessage(websocket.TextMessage, []byte("workspace updated"))
}
