package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

// LogRequestMiddleware เป็น middleware สำหรับ log ข้อมูล request ทุกชนิด
func LogRequestMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		type logData struct {
			Method  string            `json:"method"`
			Path    string            `json:"path"`
			Headers map[string]string `json:"headers"`
			Query   map[string]string `json:"query"`
			Cookies map[string]string `json:"cookies"`
			Body    string            `json:"body"`
		}

		// อ่านข้อมูล
		method := c.Method()
		path := c.OriginalURL()

		headers := make(map[string]string)
		c.Request().Header.VisitAll(func(k, v []byte) {
			headers[string(k)] = string(v)
		})

		query := c.Queries()

		cookies := make(map[string]string)
		c.Request().Header.VisitAllCookie(func(k, v []byte) {
			cookies[string(k)] = string(v)
		})

		body := c.Body()
		c.Request().SetBody(body) // คืน body ให้สามารถใช้งานต่อได้

		log := logData{
			Method:  method,
			Path:    path,
			Headers: headers,
			Query:   query,
			Cookies: cookies,
			Body:    string(body),
		}

		// log เป็น JSON
		b, _ := json.MarshalIndent(log, "", "  ")
		fmt.Println("[Request Log]:\n", string(b))

		return c.Next()
	}
}
