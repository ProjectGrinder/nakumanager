package ws

import (
	"fmt"
	"log"
	"sync"
)

type WSManager struct {
	mu          sync.Mutex
	connections map[string]map[string]WSConn
}

var Manager = WSManager{
	connections: make(map[string]map[string]WSConn),
}

func RoomKey(roomType, id string) string {
	return fmt.Sprintf("%s_%s", roomType, id)
}

func RegisterToRoom(userID string, conn WSConn, roomType, id string) {
	key := RoomKey(roomType, id)

	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	if Manager.connections[key] == nil {
		Manager.connections[key] = make(map[string]WSConn)
	}
	Manager.connections[key][userID] = conn

	log.Printf("Registered user %s to room %s", userID, key)
}

func UnregisterFromRoom(userID, roomType, id string) {
	key := RoomKey(roomType, id)

	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	if Manager.connections[key] != nil {
		delete(Manager.connections[key], userID)
		if len(Manager.connections[key]) == 0 {
			delete(Manager.connections, key)
		}
	}

	log.Printf("Unregistered user %s from room %s", userID, key)
}

func SetConnectionForTest(key, userID string, conn WSConn) {
	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	if Manager.connections == nil {
		Manager.connections = make(map[string]map[string]WSConn)
	}
	if Manager.connections[key] == nil {
		Manager.connections[key] = make(map[string]WSConn)
	}
	Manager.connections[key][userID] = conn
}

func HasRoom(key string) bool {
	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	_, exists := Manager.connections[key]
	return exists
}
