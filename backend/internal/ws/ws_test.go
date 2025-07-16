package ws_test

import (
	"sync"
	"testing"

	"github.com/nack098/nakumanager/internal/ws"
	"github.com/stretchr/testify/assert"
)

var mu sync.Mutex

func TestRegisterAndBroadcast(t *testing.T) {
	mock := &ws.MockConn{}
	userID := "u1"
	roomType := "workspace"
	roomID := "abc"

	ws.RegisterToRoom(userID, mock, roomType, roomID)
	ws.BroadcastToRoom(roomType, roomID, "test_event", "hello")

	if len(mock.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(mock.Messages))
	}
}

func TestUnregisterFromRoom(t *testing.T) {
	userID := "user-1"
	roomType := "workspace"
	roomID := "room-123"
	key := roomType + "_" + roomID

	mock := &ws.MockConn{}
	ws.SetConnectionForTest(key, userID, mock)

	ws.UnregisterFromRoom(userID, roomType, roomID)

	exists := ws.HasRoom(key)
	assert.False(t, exists, "Room should be removed after last user leaves")
}

func TestBroadcastToRoom(t *testing.T) {
	roomType := "workspace"
	roomID := "room-123"
	event := "test_event"
	payload := map[string]string{"msg": "hello"}
	key := ws.RoomKey(roomType, roomID)
	userID1 := "user1"
	userID2 := "user2"

	mock1 := &ws.MockConn{}
	mock2 := &ws.MockConn{}

	ws.SetConnectionForTest(key, userID1, mock1)
	ws.SetConnectionForTest(key, userID2, mock2)

	mu.Lock()
	ws.BroadcastToRoom(roomType, roomID, event, payload)
	defer mu.Unlock()

	assert.Len(t, mock1.Messages, 1, "mock1 should receive 1 message")
	assert.Len(t, mock2.Messages, 1, "mock2 should receive 1 message")

	msg1 := mock1.Messages[0].(map[string]interface{})
	msg2 := mock2.Messages[0].(map[string]interface{})

	assert.Equal(t, event, msg1["type"])
	assert.Equal(t, payload, msg1["data"])
	assert.Equal(t, event, msg2["type"])
	assert.Equal(t, payload, msg2["data"])
}
