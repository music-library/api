package api

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(router fiber.Router) {
	router.All("/", BaseHandler)

	// Library
	router.Get("/lib/:libId?", MainHandler)
	router.Get("/lib/:libId/tracks", TracksHandler)
	router.Get("/lib/:libId/tracks/search/:query", SearchHandler)

	router.Get("/lib/:libId/tracks/:id", TrackHandler)
	router.Get("/lib/:libId/tracks/:id/audio", TrackAudioHandler)
	router.Get("/lib/:libId/tracks/:id/cover/:size?", TrackCoverHandler)

	// Health
	router.Get("/health", HealthHandler)
	router.Get("/health/metrics", HealthHandler) // Prometheus style metrics?
}
