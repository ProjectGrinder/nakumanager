package auth_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nack098/nakumanager/internal/auth"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	app := fiber.New()
	api := app.Group("/api")
	auth.SetUpAuthRoutes(api)
	return app
}

func TestAuthRequired_MissingToken(t *testing.T) {
	app := fiber.New()
	app.Use(auth.AuthRequired)
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("OK") })

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Missing authentication token")
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(auth.AuthRequired)
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("OK") })

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid.token.value"})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid or expired token")
}

func TestAuthRequired_ValidToken(t *testing.T) {
	app := fiber.New()
	app.Use(auth.AuthRequired)
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("OK") })

	user := models.User{ID: "uid123", Email: "test@example.com", Username: "test"}
	token, err := auth.CreateToken(user)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "OK", string(body))
}

func TestRegisterAndLogin_Success(t *testing.T) {
	app := setupApp()
	registerPayload := `{
		"username": "testuser",
		"email": "test@example.com",
		"password": "VeryStrongPass123!@#"
	}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(registerPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Register Response:", string(bodyBytes))
	loginPayload := `{
		"email": "test@example.com",
		"password": "VeryStrongPass123!@#"
	}`
	req = httptest.NewRequest("POST", "/api/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, _ = io.ReadAll(resp.Body)
	fmt.Println("Login Response:", string(bodyBytes))

	var result map[string]string
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", result["message"])

	cookies := resp.Cookies()
	foundToken := false
	for _, cookie := range cookies {
		if cookie.Name == "token" && cookie.Value != "" {
			foundToken = true
			break
		}
	}
	assert.True(t, foundToken, "Token cookie should be set after login")
}

func TestRegister_InvalidEmailFormat(t *testing.T) {
	app := setupApp()

	payload := `{
		"username": "testuser",
		"email": "not-an-email",
		"password": "adjsdasdaj"
	}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestRegister_WeakPassword(t *testing.T) {
	app := setupApp()

	payload := `{
		"username": "weakuser",
		"email": "weak@example.com",
		"password": "123"
	}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestLogin_InvalidEmailFormat(t *testing.T) {
	app := setupApp()

	payload := `{
		"email": "not-an-email",
		"password": "any"
	}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestLogin_FailedAttemptsRateLimit(t *testing.T) {
	app := setupApp()

	registerPayload := `{
		"username": "limituser",
		"email": "limit@example.com",
		"password": "StrongRateLimit!123"
	}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(registerPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	for i := 0; i < 6; i++ {
		loginPayload := `{
			"email": "limit@example.com",
			"password": "WrongPassword"
		}`
		req := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginPayload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)

		if i < 5 {
			assert.Equal(t, 401, resp.StatusCode, "Attempt %d should be unauthorized", i+1)
		} else {
			assert.Equal(t, 429, resp.StatusCode, "Attempt %d should be rate-limited", i+1)
		}
	}
}

func TestLogin_InvalidEmailNotExists(t *testing.T) {
	app := setupApp()

	auth.LoginLock.Lock()
	auth.LoginAttempts = make(map[string]int)
	auth.LastAttempt = make(map[string]time.Time)
	auth.LoginLock.Unlock()

	payload := `{
		"email": "nonexistent@example.com",
		"password": "anyPassword123"
	}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "Invalid email or password")
}

func TestResetLoginAttempts(t *testing.T) {
	ip := "1.2.3.4"
	auth.RateLimitMax = 5

	auth.LoginLock.Lock()
	auth.LoginAttempts[ip] = 3
	auth.LastAttempt[ip] = time.Now()
	auth.LoginLock.Unlock()

	auth.RateLimitWindow = 100 * time.Millisecond

	auth.ResetLoginAttempts(ip)
	time.Sleep(200 * time.Millisecond)

	auth.LoginLock.Lock()
	defer auth.LoginLock.Unlock()
	assert.Equal(t, 0, auth.LoginAttempts[ip])
	_, exists := auth.LastAttempt[ip]
	assert.False(t, exists)
}

func TestVerifyToken_Invalid(t *testing.T) {
	_, err := auth.VerifyToken("invalid.token.string")
	assert.Error(t, err)
	user := models.User{ID: "uid123"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte("secret-key"))

	_, err = auth.VerifyToken(tokenStr)
	assert.Error(t, err)
}

func TestVerifyToken_Expired(t *testing.T) {
	user := models.User{ID: "uid123"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret-key"))

	_, err := auth.VerifyToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, fiber.ErrUnauthorized, err)
}

func TestBodyParser_InvalidJSON_Register(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(`{invalid-json}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "invalid character")
}

func TestBodyParser_InvalidJSON_Login(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{invalid-json}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid request body")
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	app := setupApp()

	payload := `{"username":"existuser","email":"exist@example.com","password":"VeryStrongPass123!@#"}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	req2 := httptest.NewRequest("POST", "/api/register", strings.NewReader(payload))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp2.StatusCode)

	body, _ := io.ReadAll(resp2.Body)
	assert.Contains(t, string(body), "User already exists")
}

func TestLogin_InvalidPassword(t *testing.T) {
	app := setupApp()

	payload := `{"username":"loginuser","email":"loginuser@example.com","password":"VeryStrongPass123!@#"}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	loginPayload := `{"email":"loginuser@example.com","password":"WrongPass"}`
	req2 := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginPayload))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp2.StatusCode)

	body, _ := io.ReadAll(resp2.Body)
	assert.Contains(t, string(body), "Invalid email or password")
}

func TestLogin_ValidPassword_SetsTokenCookie(t *testing.T) {
	app := setupApp()

	payload := `{"username":"tokenuser","email":"tokenuser@example.com","password":"VeryStrongPass123!@#"}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	loginPayload := `{"email":"tokenuser@example.com","password":"VeryStrongPass123!@#"}`
	req2 := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginPayload))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp2.StatusCode)

	cookies := resp2.Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "token" && c.Value != "" {
			found = true
			break
		}
	}
	assert.True(t, found, "token cookie should be set")
}
func TestAuthRequired_MissingUserIDClaim(t *testing.T) {
	app := fiber.New()
	app.Use(auth.AuthRequired)
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("OK") })

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("secret-key"))
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenStr})
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid user ID in token")
}
