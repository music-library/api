package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/config"
	"gitlab.com/music-library/music-api/version"
)

type BaseRes struct {
	Message string   `json:"message"`
	Version string   `json:"version"`
	Uptime  string   `json:"uptime"`
	Routes  []string `json:"routes"`
}

func BaseHandler(c *fiber.Ctx) error {
	return c.JSON(BaseRes{
		Message: "Hello, World 👋!",
		Version: version.Version,
		Uptime:  time.Since(config.Config.ServerStartTime).String(),
		Routes: []string{
			"/",
			"/main",
			"/tracks",
			"/track/:id",
			"/track/:id/audio",
			"/track/:id/cover/:size?",
			"/health",
			"/health/metrics",
		},
	})
}
