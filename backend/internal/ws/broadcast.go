package ws

import (
	"fmt"
	"log"
)

func BroadcastToRoom(roomType, id, event string, payload interface{}) {
	log.Printf("Broadcasting event %s to %s:%s", event, roomType, id)
	key := RoomKey(roomType, id)

	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	conns, exists := Manager.connections[key]
	if !exists {
		return
	}

	for userID, conn := range conns {
		if conn == nil {
			continue
		}
		err := conn.WriteJSON(map[string]interface{}{
			"type": event,
			"data": payload,
		})
		if err != nil {
			fmt.Printf("websocket send to user %s failed: %v\n", userID, err)
			conn.Close()
			delete(conns, userID)
		}
	}

	log.Printf("Broadcasting event %s to room %s: %d users", event, key, len(conns))

}
