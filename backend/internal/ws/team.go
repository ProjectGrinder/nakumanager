package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

func UpdateTeamHandler(c *websocket.Conn, data json.RawMessage) {
	c.WriteMessage(websocket.TextMessage, []byte("team updated"))

}
