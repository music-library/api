package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/global"
)

func ApiRoutes(router fiber.Router) {
	router.All("/", BaseHandler)

	// Library
	router.Get("/lib/:libId?", LibIdPatchMiddleware, MainHandler)
	router.Get("/lib/:libId/tracks", LibIdPatchMiddleware, TracksHandler)
	router.Get("/lib/:libId/tracks/search/:query", LibIdPatchMiddleware, SearchHandler)

	router.Get("/lib/:libId/tracks/:id", LibIdPatchMiddleware, TrackHandler)
	router.Get("/lib/:libId/tracks/:id/audio", LibIdPatchMiddleware, TrackAudioHandler)
	router.Get("/lib/:libId/tracks/:id/cover/:size?", LibIdPatchMiddleware, TrackCoverHandler)

	// Health
	router.Get("/health", HealthHandler)
	router.Get("/health/metrics", HealthHandler) // Prometheus style metrics?
}

func LibIdPatchMiddleware(c *fiber.Ctx) error {
	libId := c.Params("libId", global.IndexMany.DefaultKey)

	// @Node: Header has a problem: We can't set one on the FE for things like img src or audio src
	// libId := c.Get("X-Library", global.IndexMany.DefaultKey)

	c.Locals("libId", libId)

	if _, ok := global.IndexMany.Indexes[libId]; !ok {
		return Error(c, 404, fmt.Sprintf("library `%s` does not exist", libId))
	}

	return c.Next()
}
