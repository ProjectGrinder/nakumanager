package ws

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nack098/nakumanager/internal/repositories"
)

type ConnWithLocals interface {
	Locals(key string, defaultValue ...interface{}) interface{}
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
}

type WSHandler struct {
	mu            sync.Mutex
	Clients       map[*websocket.Conn]string
	WorkspaceRepo repositories.WorkspaceRepository
	TeamRepo      repositories.TeamRepository
	ProjectRepo   repositories.ProjectRepository
	IssueRepo     repositories.IssueRepository
	UserRepo      repositories.UserRepository
	ViewRepo      repositories.ViewRepository
	BroadcastFunc func(interface{})
}

func NewWSHandler(workspaceRepo repositories.WorkspaceRepository, teamRepo repositories.TeamRepository, projectRepo repositories.ProjectRepository, issueRepo repositories.IssueRepository, userRepo repositories.UserRepository, viewRepo repositories.ViewRepository) *WSHandler {
	return &WSHandler{
		Clients:       make(map[*websocket.Conn]string),
		WorkspaceRepo: workspaceRepo,
		TeamRepo:      teamRepo,
		ProjectRepo:   projectRepo,
		IssueRepo:     issueRepo,
		UserRepo:      userRepo,
		ViewRepo:      viewRepo,
	}
}

func (h *WSHandler) RegisterClient(conn *websocket.Conn, client string) error {
	if conn == nil || client == "" {
		return errors.New("invalid client registration: nil connection or empty client ID")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.Clients[conn]; ok {
		return errors.New("client already registered")
	}
	h.Clients[conn] = client
	log.Printf("New WebSocket client connected. Total: %d", len(h.Clients))
	return nil
}

func (h *WSHandler) UnregisterClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.Clients, conn)
	log.Printf("Client disconnected. Remaining: %d", len(h.Clients))
}

func (h *WSHandler) Broadcast(message interface{}) {
	if h.BroadcastFunc != nil {
		h.BroadcastFunc(message)
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("broadcast marshal error:", err)
		return
	}
	for conn := range h.Clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write message error:", err)
		}
	}
}

func WebSocketMiddleware(authHandler interface {
	VerifyToken(string) (*jwt.Token, error)
}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		token := c.Cookies("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authentication token",
			})
		}

		_, err := authHandler.VerifyToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		claims, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		})
		if m, ok := claims.Claims.(jwt.MapClaims); ok {
			if userID, ok := m["user_id"].(string); ok {
				c.Locals("userID", userID)
			}
		}

		return c.Next()
	}
}

func (h *WSHandler) CentralWebSocketHandler(c *websocket.Conn) {
	defer func() {
		h.UnregisterClient(c)
		c.Close()
	}()

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		log.Println("Missing userID in connection")
		return
	}

	h.RegisterClient(c, userID)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var message struct {
			Event string          `json:"event"`
			Data  json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("json error:", err)
			continue
		}

		log.Println("Received event:", message.Event)
		switch message.Event {
		case "update_workspace":
			h.UpdateWorkspaceHandler(c, message.Data)
		case "update_project":
			h.UpdateProjectHandler(c, message.Data)
		case "update_issue":
			h.UpdateIssueHandler(c, message.Data)
		case "update_view":
			h.UpdateViewHandler(c, message.Data)
		case "update_team":
			h.UpdateTeamHandler(c, message.Data)
		default:
			log.Println("unknown event:", message.Event)
		}
	}
}
