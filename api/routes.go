package api

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(router fiber.Router) {
	// Track
	router.Get("/tracks", TracksHandler)
	router.Get("/track/:id", TrackHandler)
	router.Get("/track/:id/audio", TrackAudioHandler)
	router.Get("/track/:id/cover/:size?", TrackCoverHandler)

	// Album
	router.Get("/albums", MockHandler)

	// Health
	router.Get("/health", MockHandler)
	router.Get("/health/metrics", MockHandler) // Prometheus style metrics?
}
