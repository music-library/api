package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/indexer"
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
	libId := c.Params("libId", indexer.MusicLibIndex.DefaultKey)

	// @Node: Header has a problem: We can't set one on the FE for things like img src or audio src
	// libId := c.Get("X-Library", indexer.MusicLibIndex.DefaultKey)

	c.Locals("libId", libId)

	if _, ok := indexer.MusicLibIndex.Indexes[libId]; !ok {
		return Error(c, 404, fmt.Sprintf("library `%s` does not exist", libId))
	}

	return c.Next()
}
