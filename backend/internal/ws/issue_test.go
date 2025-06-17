package ws_test

import (
	"encoding/json"
	"net/url"
	"testing"

	fiberWs "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gorillaWs "github.com/gorilla/websocket"
	"github.com/nack098/nakumanager/internal/ws"
	"github.com/stretchr/testify/assert"
)

func setupWebSocketIssueApp() *fiber.App {
	app := fiber.New()
	wsGroup := app.Group("/ws", ws.WebSocketMiddleware)
	wsGroup.Get("/", fiberWs.New(ws.CentralWebSocketHandler))
	return app
}

func TestUpdateIssueHandler(t *testing.T) {
	app := setupWebSocketIssueApp()

	go func() {
		err := app.Listen(":3000")
		if err != nil {
			t.Log("Server error:", err)
		}
	}()
	defer app.Shutdown()

	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/ws/"}

	conn, _, err := gorillaWs.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	message := map[string]interface{}{
		"event": "update_issue",
		"data":  map[string]interface{}{},
	}
	payload, _ := json.Marshal(message)
	err = conn.WriteMessage(gorillaWs.TextMessage, payload)
	assert.NoError(t, err)

	_, resp, err := conn.ReadMessage()
	assert.NoError(t, err)

	assert.Equal(t, "issue updated", string(resp))
}
