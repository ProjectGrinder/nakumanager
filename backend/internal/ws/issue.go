package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

func UpdateIssueHandler(c *websocket.Conn, data json.RawMessage) {
	c.WriteMessage(websocket.TextMessage, []byte("issue updated"))

}
