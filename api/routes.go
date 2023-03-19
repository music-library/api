package api

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(router fiber.Router) {
	router.All("/", BaseHandler)

	router.Get("/main", MainHandler)

	// Track
	router.Get("/tracks", TracksHandler)
	router.Get("/tracks/search/:query", SearchHandler)

	router.Get("/tracks/:id", TrackHandler)
	router.Get("/tracks/:id/audio", TrackAudioHandler)
	router.Get("/tracks/:id/cover/:size?", TrackCoverHandler)

	// Library
	router.Get("/lib/:libId/tracks", TracksHandler)

	// Health
	router.Get("/health", HealthHandler)
	router.Get("/health/metrics", HealthHandler) // Prometheus style metrics?
}
