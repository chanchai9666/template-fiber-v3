package usecases

import (
	"github.com/gofiber/fiber/v3"
)

func GetUserID(c fiber.Ctx) string {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return ""
	}
	return userID
}
