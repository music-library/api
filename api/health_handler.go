package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/global"
)

func HealthHandler(c *fiber.Ctx) error {
	message := "ok"
	status := 200
	ok := true

	for _, index := range global.IndexMany.Indexes {
		if len(index.Tracks) == 0 {
			message = fmt.Sprintf("%s track index is empty", index.Name)
			status = 500
			ok = false
			break
		}
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"ok":      ok,
	})
}
