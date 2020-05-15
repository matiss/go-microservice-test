package server

import (
	"github.com/gofiber/fiber"
)

type ErrorResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func RootHandler(c *fiber.Ctx) {
	c.JSON(
		map[string]interface{}{
			"message": "Method Not Allowed",
		},
	)
}

func RobotsTXTHandler(c *fiber.Ctx) {
	c.Send("User-agent: *\nDisallow: /")
}
