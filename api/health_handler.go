package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/global"
)

func HealthHandler(c *fiber.Ctx) error {
	message := "ok"
	status := 200
	ok := true

	if global.Index.TracksCount == 0 {
		message = "track index is empty"
		status = 500
		ok = false
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"ok":      ok,
	})
}
