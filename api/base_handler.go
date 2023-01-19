package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/version"
)

func BaseHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World 👋!",
		"version": version.Version,
		"routes": []string{
			"/",
			"/main",
			"/tracks",
			"/track/:id",
			"/track/:id/audio",
			"/track/:id/cover/:size?",
			"/albums",
			"/health",
			"/health/metrics",
		},
	})
}
