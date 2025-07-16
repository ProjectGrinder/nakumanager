package auth_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wsfiber "github.com/gofiber/contrib/websocket"
	"github.com/nack098/nakumanager/internal/auth"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, data db.CreateUserParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id string) (db.GetUserByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.GetUserByIDRow), args.Error(1)
}

func (m *MockUserRepo) ListUsers(ctx context.Context) ([]db.ListUsersRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.ListUsersRow), args.Error(1)
}

func (m *MockUserRepo) UpdateEmail(ctx context.Context, data db.UpdateEmailParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateRoles(ctx context.Context, data db.UpdateRolesParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateUsername(ctx context.Context, data db.UpdateUsernameParams) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailWithoutPasswordRow, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(db.GetUserByEmailWithoutPasswordRow), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmailWithPassword(ctx context.Context, email string) (db.GetUserByEmailWithPasswordRow, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(db.GetUserByEmailWithPasswordRow), args.Error(1)
}

func setupApp(handler *auth.AuthHandler) *fiber.App {
	app := fiber.New()
	app.Post("/login", handler.Login)
	app.Post("/register", handler.Register)
	return app
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)

	email := "test@example.com"
	password := "VeryStrongPassword!@123##$@"

	mockRepo.On("GetUserByEmail", mock.Anything, email).
		Return(db.GetUserByEmailWithoutPasswordRow{}, errors.New("not found"))

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).
		Return(nil)

	reqBody := fmt.Sprintf(`{
		"username": "tester",
		"email": "%s",
		"password": "%s"
	}`, email, password)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "User registered successfully")

	mockRepo.AssertExpectations(t)
}

func TestRegister_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)

	reqBody := `{
		"username": "tester",
		"email": "not-an-email",
		"password": "Strong123!"
	}`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid email format")
}

func TestRegister_BodyParserError(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)

	badJSON := `{"username": "tester", "email": "test@example.com", "password": "abc"`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(badJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)

	assert.Contains(t, string(body), "unexpected end of JSON input")
}

func TestRegister_WeakPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)

	email := "test@example.com"
	reqBody := fmt.Sprintf(`{
		"username": "tester",
		"email": "%s",
		"password": "abc"
	}`, email)

	mockRepo.On("GetUserByEmail", mock.Anything, email).
		Return(db.GetUserByEmailWithoutPasswordRow{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Password is too weak")
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)

	email := "test@example.com"
	mockRepo.On("GetUserByEmail", mock.Anything, email).
		Return(db.GetUserByEmailWithoutPasswordRow{}, nil)

	reqBody := fmt.Sprintf(`{
		"username": "tester",
		"email": "%s",
		"password": "Strong123!!!@#$"
	}`, email)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "User already exists")
}

func TestRegister_CreateUserFail(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)

	email := "test@example.com"
	password := "VeryStrongPassword!@123##$@"

	mockRepo.On("GetUserByEmail", mock.Anything, email).
		Return(db.GetUserByEmailWithoutPasswordRow{}, errors.New("not found"))

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).
		Return(errors.New("db error"))

	reqBody := fmt.Sprintf(`{
		"username": "tester",
		"email": "%s",
		"password": "%s"
	}`, email, password)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Failed to create user")
}

func resetLoginRateLimit() {
	auth.LoginLock.Lock()
	auth.LoginAttempts = make(map[string]int)
	auth.LastAttempt = make(map[string]time.Time)
	auth.LoginLock.Unlock()
}

func TestLogin_BodyParserError(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)
	resetLoginRateLimit()

	badJSON := `{"email": "test@example.com", "password": "password123"`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(badJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid request body")

	mockRepo.AssertNotCalled(t, "GetUserByEmailWithPassword", mock.Anything, mock.Anything)
}

func TestLogin_RateLimitExceeded(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)
	resetLoginRateLimit()

	ip := "0.0.0.0"

	auth.LoginLock.Lock()
	auth.LoginAttempts = map[string]int{
		ip: auth.RateLimitMax,
	}
	auth.LastAttempt = make(map[string]time.Time)
	auth.LoginLock.Unlock()

	mockRepo.On("GetUserByEmailWithPassword", mock.Anything, mock.Anything).
		Return(db.GetUserByEmailWithPasswordRow{
			ID:           "user-id-123",
			Username:     "tester",
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
			Roles:        "user",
		}, nil)

	reqBody := `{"email":"test@example.com","password":"StrongPassword!123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("X-Forwarded-For", ip)
	req.RemoteAddr = ip + ":1234"

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Too many login attempts")

	mockRepo.AssertNotCalled(t, "GetUserByEmailWithPassword", mock.Anything, mock.Anything)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)
	resetLoginRateLimit()

	email := "notfound@example.com"
	password := "AnyPassword123!"

	mockRepo.On("GetUserByEmailWithPassword", mock.Anything, email).
		Return(db.GetUserByEmailWithPasswordRow{}, errors.New("not found"))

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid email or password")

	mockRepo.AssertExpectations(t)
}
func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)

	email := "user@example.com"
	wrongPassword := "wrongpassword"
	correctPassword := "correctpassword"

	hashPass, err := argon2id.CreateHash(correctPassword, argon2id.DefaultParams)
	require.NoError(t, err)

	mockRepo.On("GetUserByEmailWithPassword", mock.Anything, email).
		Return(db.GetUserByEmailWithPasswordRow{
			ID:           "user-id-123",
			Username:     "tester",
			Email:        email,
			PasswordHash: hashPass,
			Roles:        "user",
		}, nil)

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, wrongPassword)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid email or password")

	mockRepo.AssertExpectations(t)
}

func TestLogin_ValidPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)

	email := "user@example.com"
	password := "correctpassword"

	hashPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	require.NoError(t, err)

	mockRepo.On("GetUserByEmailWithPassword", mock.Anything, email).
		Return(db.GetUserByEmailWithPasswordRow{
			ID:           "user-id-123",
			Username:     "tester",
			Email:        email,
			PasswordHash: hashPass,
			Roles:        "user",
		}, nil)

	handler.CreateTokenFunc = func(user models.User) (string, error) {
		return "valid.jwt.token", nil
	}

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Login successful")

	mockRepo.AssertExpectations(t)
}

func TestLogin_Success_WithRealArgon2id(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)
	resetLoginRateLimit()

	password := "StrongPassword123!"
	hashPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	assert.NoError(t, err)

	mockRepo.On("GetUserByEmailWithPassword", mock.Anything, "test@example.com").
		Return(db.GetUserByEmailWithPasswordRow{
			ID:           "user-id-123",
			Username:     "tester",
			Email:        "test@example.com",
			PasswordHash: hashPass,
			Roles:        "user",
		}, nil)

	reqBody := fmt.Sprintf(`{"email":"test@example.com","password":"%s"}`, password)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Login successful")

	mockRepo.AssertExpectations(t)
}

func TestLogin_CreateTokenFail(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	resetLoginRateLimit()

	app := setupApp(handler)

	password := "StrongPassword123!"
	hashPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	assert.NoError(t, err)

	email := "test@example.com"

	mockRepo.On("GetUserByEmailWithPassword", mock.Anything, email).
		Return(db.GetUserByEmailWithPasswordRow{
			ID:           "user-id-123",
			Username:     "tester",
			Email:        email,
			PasswordHash: hashPass,
			Roles:        "user",
		}, nil)

	handler.CreateTokenFunc = func(user models.User) (string, error) {
		return "", errors.New("token creation failed")
	}

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Error while creating token")

	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmailFormat(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)
	resetLoginRateLimit()

	badEmail := "invalid-email-format"
	reqBody := fmt.Sprintf(`{"email":"%s","password":"AnyPassword123"}`, badEmail)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid email format")

	mockRepo.AssertNotCalled(t, "GetUserByEmailWithPassword", mock.Anything, mock.Anything)
}

func TestLogin_Success_UsingRealCompare(t *testing.T) {
	mockRepo := new(MockUserRepo)
	handler := auth.NewAuthHandler(mockRepo)
	app := setupApp(handler)
	resetLoginRateLimit()

	email := "test@example.com"
	password := "VeryStrongPassword!@#"

	hashPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	require.NoError(t, err)

	mockRepo.On("GetUserByEmailWithPassword", mock.Anything, email).
		Return(db.GetUserByEmailWithPasswordRow{
			ID:           "user-id-123",
			PasswordHash: hashPass,
		}, nil)

	reqBody := fmt.Sprintf(`{
		"email": "%s",
		"password": "%s"
	}`, email, password)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "Login successful")

	mockRepo.AssertExpectations(t)
}

func TestResetLoginAttempts(t *testing.T) {
	ip := "192.168.0.1"

	oldWindow := auth.RateLimitWindow
	auth.RateLimitWindow = 10 * time.Millisecond
	defer func() { auth.RateLimitWindow = oldWindow }()

	auth.LoginLock.Lock()
	auth.LoginAttempts = map[string]int{ip: 5}
	auth.LastAttempt = map[string]time.Time{ip: time.Now()}
	auth.LoginLock.Unlock()

	auth.ResetLoginAttempts(ip)

	time.Sleep(auth.RateLimitWindow + 20*time.Millisecond)

	auth.LoginLock.Lock()
	defer auth.LoginLock.Unlock()

	if attempts := auth.LoginAttempts[ip]; attempts != 0 {
		t.Errorf("expected LoginAttempts[%s] to be 0, got %d", ip, attempts)
	}

	if _, exists := auth.LastAttempt[ip]; exists {
		t.Errorf("expected LastAttempt[%s] to be deleted, but still exists", ip)
	}
}

func setupAuthRequiredTestApp(handler *auth.AuthHandler) *fiber.App {
	app := fiber.New()
	app.Use(handler.AuthRequired)
	app.Get("/protected", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		return c.JSON(fiber.Map{"userID": userID})
	})
	return app
}

func TestAuthRequired_NoTokenCookie(t *testing.T) {
	handler := &auth.AuthHandler{}
	app := setupAuthRequiredTestApp(handler)

	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Missing authentication token")
}

func TestAuthRequired_VerifyTokenError(t *testing.T) {
	handler := &auth.AuthHandler{
		VerifyTokenFunc: func(tokenStr string) (*jwt.Token, error) {
			return nil, errors.New("invalid token")
		},
	}
	app := setupAuthRequiredTestApp(handler)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "bad.token"})
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid or expired token")
}

func TestAuthRequired_MissingOrEmptyUserID(t *testing.T) {
	handler := &auth.AuthHandler{
		VerifyTokenFunc: func(tokenStr string) (*jwt.Token, error) {
			return &jwt.Token{
				Claims: jwt.MapClaims{},
				Valid:  true,
			}, nil
		},
	}
	app := setupAuthRequiredTestApp(handler)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "token.without.userid"})
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid user ID in token")
}

func TestAuthRequired_Success(t *testing.T) {
	handler := &auth.AuthHandler{
		VerifyTokenFunc: func(tokenStr string) (*jwt.Token, error) {
			return &jwt.Token{
				Claims: jwt.MapClaims{
					"user_id": "12345",
				},
				Valid: true,
			}, nil
		},
	}
	app := setupAuthRequiredTestApp(handler)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "valid.token"})
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "12345")
}

var secretKey = []byte("secret-key")

func TestVerifyToken_ValidToken(t *testing.T) {
	handler := &auth.AuthHandler{}

	claims := jwt.MapClaims{
		"user_id": "12345",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	require.NoError(t, err)

	parsedToken, err := handler.VerifyToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)
	assert.True(t, parsedToken.Valid)

	mapClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, "12345", mapClaims["user_id"])
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	handler := &auth.AuthHandler{}

	invalidTokenString := "this.is.an.invalid.token"

	parsedToken, err := handler.VerifyToken(invalidTokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedToken)
	assert.Equal(t, fiber.ErrUnauthorized, err)
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	handler := &auth.AuthHandler{}

	claims := jwt.MapClaims{
		"user_id": "12345",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	require.NoError(t, err)

	parsedToken, err := handler.VerifyToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedToken)
	assert.Equal(t, fiber.ErrUnauthorized, err)
}

func setupWebsocketTestApp(handler *auth.AuthHandler) *fiber.App {
	app := fiber.New()
	app.Use(handler.WebSocketAuthRequired())

	app.Get("/ws", wsfiber.New(func(c *wsfiber.Conn) {
		userID := c.Locals("userID").(string)
		_ = c.WriteMessage(websocket.TextMessage, []byte("userID:"+userID))
	}))

	return app
}

func TestWebSocketAuthRequired_Success(t *testing.T) {
	app := fiber.New()

	authHandler := &auth.AuthHandler{
		VerifyTokenFunc: func(tokenStr string) (*jwt.Token, error) {
			return &jwt.Token{
				Claims: jwt.MapClaims{
					"user_id": "12345",
				},
				Valid: true,
			}, nil
		},
	}

	app.Use("/ws", authHandler.WebSocketAuthRequired())

	app.Get("/ws", wsfiber.New(func(conn *wsfiber.Conn) {
		userID := conn.Locals("userID").(string)
		_ = conn.WriteMessage(websocket.TextMessage, []byte("userID:"+userID))
	}))

	go func() {
		_ = app.Listen(":3001")
	}()
	defer app.Shutdown()

	time.Sleep(100 * time.Millisecond) 

	dialer := websocket.Dialer{}
	header := http.Header{}
	header.Add("Cookie", "token=valid.token")

	wsURL := "ws://localhost:3001/ws"
	conn, resp, err := dialer.Dial(wsURL, header)

	require.NoError(t, err)
	require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, "userID:12345", string(msg))

	_ = conn.Close()
}

func TestWebSocketAuth_MissingToken(t *testing.T) {
	handler := &auth.AuthHandler{}
	app := fiber.New()

	app.Use(handler.WebSocketAuthRequired())
	app.Get("/ws", func(c *fiber.Ctx) error {
		return c.SendString("Should not reach here")
	})

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Missing token")
}

func TestWebSocketAuth_InvalidToken(t *testing.T) {
	handler := &auth.AuthHandler{
		VerifyTokenFunc: func(tokenStr string) (*jwt.Token, error) {
			return nil, errors.New("invalid token")
		},
	}
	app := fiber.New()

	app.Use(handler.WebSocketAuthRequired())
	app.Get("/ws", func(c *fiber.Ctx) error {
		return c.SendString("Should not reach here")
	})

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid.token"})
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid or expired token")
}

func TestWebSocketAuth_MissingUserID(t *testing.T) {
	handler := &auth.AuthHandler{
		VerifyTokenFunc: func(tokenStr string) (*jwt.Token, error) {
			claims := jwt.MapClaims{} 
			return &jwt.Token{
				Claims: claims,
				Valid:  true,
			}, nil
		},
	}
	app := fiber.New()

	app.Use(handler.WebSocketAuthRequired())
	app.Get("/ws", func(c *fiber.Ctx) error {
		return c.SendString("Should not reach here")
	})

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "valid.token"})
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid user ID")
}
