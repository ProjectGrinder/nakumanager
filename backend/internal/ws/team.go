package ws

import(
	"log"
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
)


func UpdateTeamHandler(c *websocket.Conn, data json.RawMessage) {
	log.Println("update team")

}