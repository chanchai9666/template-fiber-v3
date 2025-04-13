package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v3"

	"template-fiber-v3/internal/pkg/safety"
)

// AuthMiddleware เป็น middleware สำหรับตรวจสอบ JWT และ Basic Auth
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Please Login before use."})
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			// JWT Authentication
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := safety.VerifyJWT(jwtSecret, tokenString)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT"})
			}

			// เก็บ claims ใน context
			c.Locals("claims", claims)
			c.Locals("userID", claims.UserId)

		} else if strings.HasPrefix(authHeader, "Basic ") {
			// Basic Authentication
			encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
			decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Basic Auth"})
			}

			credentials := strings.SplitN(string(decodedCredentials), ":", 2)
			if len(credentials) != 2 {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Basic Auth"})
			}

			username, password := credentials[0], credentials[1]

			if !validateBasicAuth(username, password) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
			}

			// ถ้าผ่าน Basic Auth อาจจะตั้ง userID เป็นชื่อผู้ใช้
			c.Locals("userID", username)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unsupported Authorization method"})
		}

		return c.Next()
	}
}

// validateBasicAuth เป็นฟังก์ชันตัวอย่างสำหรับตรวจสอบ username และ password
func validateBasicAuth(username, password string) bool {
	return username == "admin" && password == "admin"
}

// GetUserProfile ดึง userID จาก context
func GetUserProfile(c fiber.Ctx) string {
	userID := c.Locals("userID")
	if userID == nil {
		return ""
	}
	return userID.(string)
}
