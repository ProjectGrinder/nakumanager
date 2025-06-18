package ws_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	fiberWs "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gorillaWs "github.com/gorilla/websocket"
	"github.com/nack098/nakumanager/internal/ws"
	"github.com/stretchr/testify/assert"
)

func ParseWebSocketMessage(msg []byte) (string, json.RawMessage, error) {
	var message struct {
		Event string          `json:"event"`
		Data  json.RawMessage `json:"data"`
	}
	err := json.Unmarshal(msg, &message)
	return message.Event, message.Data, err
}

func TestWebSocketMiddleware_Upgrade(t *testing.T) {
	app := fiber.New()

	app.Use(ws.WebSocketMiddleware, func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestWebSocketMiddleware_NonUpgrade(t *testing.T) {
	app := fiber.New()

	app.Use(ws.WebSocketMiddleware)

	req := httptest.NewRequest("GET", "/", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUpgradeRequired, resp.StatusCode)
}

func TestParseWebSocketMessage_Success(t *testing.T) {
	msg := []byte(`{"event": "update_project", "data": {"id": 1}}`)

	event, data, err := ParseWebSocketMessage(msg)

	assert.NoError(t, err)
	assert.Equal(t, "update_project", event)
	assert.JSONEq(t, `{"id": 1}`, string(data))
}

func TestParseWebSocketMessage_InvalidJSON(t *testing.T) {
	msg := []byte(`invalid json`)

	_, _, err := ParseWebSocketMessage(msg)

	assert.Error(t, err)
}

func setupWebSocketApp() *fiber.App {
	app := fiber.New()

	app.Use("/ws", ws.WebSocketMiddleware)
	app.Get("/ws", fiberWs.New(ws.CentralWebSocketHandler))

	return app
}

func TestCentralWebSocketHandler_JSONError(t *testing.T) {
	app := setupWebSocketApp()

	go func() {
		err := app.Listen(":3000")
		if err != nil {
			t.Log("Server error:", err)
		}
	}()
	defer app.Shutdown()

	time.Sleep(100 * time.Millisecond)

	url := "ws://localhost:3000/ws"
	wsConn, _, err := gorillaWs.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer wsConn.Close()

	err = wsConn.WriteMessage(gorillaWs.TextMessage, []byte(`invalid json`))
	assert.NoError(t, err)
}

func TestCentralWebSocketHandler_UnknownEvent(t *testing.T) {
	app := setupWebSocketApp()

	go func() {
		err := app.Listen(":3000")
		if err != nil {
			t.Log("Server error:", err)
		}
	}()
	defer app.Shutdown()

	time.Sleep(100 * time.Millisecond)

	url := "ws://localhost:3000/ws"
	wsConn, _, err := gorillaWs.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer wsConn.Close()

	message := []byte(`{"event":"some_event", "data":{}}`)
	err = wsConn.WriteMessage(gorillaWs.TextMessage, message)
	assert.NoError(t, err)
}
