package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/global"
)

func ApiRoutes(router fiber.Router) {
	router.All("/", BaseHandler)

	// Library
	router.Get("/lib", LibIdPatchMiddleware, MainHandler)
	router.Get("/lib/tracks", LibIdPatchMiddleware, TracksHandler)
	router.Get("/lib/tracks/search/:query", LibIdPatchMiddleware, SearchHandler)

	router.Get("/lib/tracks/:id", LibIdPatchMiddleware, TrackHandler)
	router.Get("/lib/tracks/:id/audio", LibIdPatchMiddleware, TrackAudioHandler)
	router.Get("/lib/tracks/:id/cover/:size?", LibIdPatchMiddleware, TrackCoverHandler)

	// Health
	router.Get("/health", HealthHandler)
	router.Get("/health/metrics", HealthHandler) // Prometheus style metrics?
}

func LibIdPatchMiddleware(c *fiber.Ctx) error {
	libId := c.Get("X-Library", global.IndexMany.DefaultKey)

	c.Locals("libId", libId)

	if _, ok := global.IndexMany.Indexes[libId]; !ok {
		return Error(c, 404, fmt.Sprintf("library `%s` does not exist", libId))
	}

	return c.Next()
}
