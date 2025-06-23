package auth

import (
	"net/mail"
	"sync"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
	"github.com/nbutton23/zxcvbn-go"
)

type AuthHandler struct {
	UserRepo        repositories.UserRepository
	CreateTokenFunc func(user models.User) (string, error)
	VerifyTokenFunc func(tokenStr string) (*jwt.Token, error)
}

func NewAuthHandler(userRepo repositories.UserRepository) *AuthHandler {
	h := &AuthHandler{
		UserRepo: userRepo,
	}
	h.CreateTokenFunc = h.CreateToken
	h.VerifyTokenFunc = h.verifyTokenInternal  
	return h
}

func (h *AuthHandler) verifyTokenInternal(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	return token, nil
}

func (h *AuthHandler) VerifyToken(tokenStr string) (*jwt.Token, error) {
	if h.VerifyTokenFunc != nil {
		return h.VerifyTokenFunc(tokenStr)
	}

	return h.verifyTokenInternal(tokenStr)
}

var (
	dummyHash, _    = argon2id.CreateHash("dummy_password", argon2id.DefaultParams)
	LoginAttempts   = make(map[string]int)
	LoginLock       = sync.Mutex{}
	RateLimitMax    = 5
	RateLimitWindow = time.Minute * 5
	LastAttempt     = make(map[string]time.Time)
	secretKey       = []byte("secret-key")
	ResetDelay = time.Minute
)

func (h *AuthHandler) CreateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(secretKey)
}
func (h *AuthHandler) AuthRequired(c *fiber.Ctx) error {
	tokenStr := c.Cookies("token")
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authentication token"})
	}

	token, err := h.VerifyToken(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	c.Locals("userID", userID)
	return c.Next()
}

func ResetLoginAttempts(ip string) {
	time.AfterFunc(RateLimitWindow, func() {
		LoginLock.Lock()
		defer LoginLock.Unlock()
		LoginAttempts[ip] = 0
		delete(LastAttempt, ip)
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body models.Login
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		return c.Status(400).SendString("Invalid email format")
	}

	ip := c.IP()

	LoginLock.Lock()
	if count, ok := LoginAttempts[ip]; ok && count >= RateLimitMax {
		LoginLock.Unlock()
		return c.Status(429).SendString("Too many login attempts, please try again later")
	}
	LoginLock.Unlock()

	user, err := h.UserRepo.GetUserByEmailWithPassword(c.Context(), body.Email)
	if err != nil {
		_, _ = argon2id.ComparePasswordAndHash(body.Password, dummyHash)
		LoginLock.Lock()
		LoginAttempts[ip]++
		LastAttempt[ip] = time.Now()
		if LoginAttempts[ip] == 1 {
			ResetLoginAttempts(ip)
		}
		LoginLock.Unlock()
		return c.Status(401).SendString("Invalid email or password")
	}

	match, err := argon2id.ComparePasswordAndHash(body.Password, user.PasswordHash)
	if err != nil || !match {
		LoginLock.Lock()
		LoginAttempts[ip]++
		LastAttempt[ip] = time.Now()
		if LoginAttempts[ip] == 1 {
			ResetLoginAttempts(ip)
		}
		LoginLock.Unlock()
		return c.Status(401).SendString("Invalid email or password")
	}

	LoginLock.Lock()
	LoginAttempts[ip] = 0
	delete(LastAttempt, ip)
	LoginLock.Unlock()

	tokenString, err := h.CreateTokenFunc(models.User{ID: user.ID})
	if err != nil {
		return c.Status(500).SendString("Error while creating token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body models.Register
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		return c.Status(400).SendString("Invalid email format")
	}

	_, err := h.UserRepo.GetUserByEmail(c.Context(), body.Email)
	if err == nil {
		return c.Status(400).SendString("User already exists")
	}

	password := zxcvbn.PasswordStrength(body.Password, nil)
	if password.Score < 3 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password is too weak",
		})
	}

	hashPass, err := argon2id.CreateHash(body.Password, argon2id.DefaultParams)
	if err != nil {
		return c.Status(500).SendString("Failed to hash password")
	}

	user := models.User{
		ID:           uuid.New().String(),
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: hashPass,
	}

	err = h.UserRepo.CreateUser(c.Context(), db.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Roles:        "user",
	})
	if err != nil {
		return c.Status(500).SendString("Failed to create user")
	}

	return c.SendString("User registered successfully!")
}
