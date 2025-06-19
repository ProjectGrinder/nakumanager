package auth

import (
	"net/mail"
	"sync"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nbutton23/zxcvbn-go"
)

var (
	dummyHash, _    = argon2id.CreateHash("dummy_password", argon2id.DefaultParams) //dummy hash ป้องกัน timming attack
	LoginAttempts   = make(map[string]int)                                          //จำนวนการ login ที่ไม่ถูกต้อง
	LoginLock       = sync.Mutex{}                                                  // mutex lock
	RateLimitMax    = 5                                                             // จำนวนครั้งสูงสุดที่ 1 ip สามารถ login ผิดได้
	RateLimitWindow = time.Minute * 5                                               // เวลารอหลัง login ผิดเงื่อนไข RateLimitMax
	LastAttempt     = make(map[string]time.Time)                                    //เวลาล่าสุดที่ login ผิด

	// Mock DB
	Users = make(map[string]models.User)
	//Mock secret key for jwt
	secretKey = []byte("secret-key")
)

func CreateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func VerifyToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}

	return nil
}

func AuthRequired(c *fiber.Ctx) error {

	token := c.Cookies("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authentication token",
		})
	}

	if err := VerifyToken(token); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

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

func SetUpAuthRoutes(api fiber.Router) {
	api.Post("/login", Login)
	api.Post("/register", Register)
}

func Login(c *fiber.Ctx) error {
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

	user, exists := Users[body.Email]
	if !exists {
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

	tokenString, err := CreateToken(user)
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

func Register(c *fiber.Ctx) error {
	var body models.Register
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		return c.Status(400).SendString("Invalid email format")
	}

	if _, exists := Users[body.Email]; exists {
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

	Users[body.Email] = user
	return c.SendString("User registered successfully!")
}
