package api

import (
	"github.com/gofiber/fiber/v2"
)

func Error(c *fiber.Ctx, status int, msg any) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"message": msg,
	})
}
