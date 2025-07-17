package ws

import (
	"encoding/json"
	"log"

	wsfiber "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)


type WebSocketHandler struct{}

type ClientMessage struct {
	Type  string              `json:"type"`  
	Rooms map[string][]string `json:"rooms"`
}

func (h *WebSocketHandler) Handle(c *fiber.Ctx) error {
	if !wsfiber.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	userIDIfc := c.Locals("userID")
	if userIDIfc == nil {
		return fiber.ErrUnauthorized
	}
	userID := userIDIfc.(string)

	return wsfiber.New(func(conn *wsfiber.Conn) {
		defer conn.Close()

		activeRooms := map[string]map[string]bool{}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("websocket disconnected: %v", err)
				break
			}

			var clientMsg ClientMessage
			if err := json.Unmarshal(msg, &clientMsg); err != nil {
				log.Println("invalid websocket message:", err)
				continue
			}

			switch clientMsg.Type {
			case "subscribe":
				for roomType, ids := range clientMsg.Rooms {
					for _, id := range ids {
						RegisterToRoom(userID, conn, roomType, id)

						
						if activeRooms[roomType] == nil {
							activeRooms[roomType] = make(map[string]bool)
						}
						activeRooms[roomType][id] = true
					}
				}
			case "unsubscribe":
				for roomType, ids := range clientMsg.Rooms {
					for _, id := range ids {
						UnregisterFromRoom(userID, roomType, id)

						
						if activeRooms[roomType] != nil {
							delete(activeRooms[roomType], id)
							if len(activeRooms[roomType]) == 0 {
								delete(activeRooms, roomType)
							}
						}
					}
				}
			default:
				log.Printf("unknown message type: %s", clientMsg.Type)
			}
		}

		
		for roomType, ids := range activeRooms {
			for id := range ids {
				UnregisterFromRoom(userID, roomType, id)
			}
		}
	})(c)
}
