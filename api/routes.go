package api

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(router fiber.Router) {
	router.All("/", BaseHandler)

	// Track
	router.Get("/tracks", TracksHandler)
	router.Get("/track/:id", TrackHandler)
	router.Get("/track/:id/audio", TrackAudioHandler)
	router.Get("/track/:id/cover/:size?", TrackCoverHandler)

	// Album
	router.Get("/albums", BaseHandler)

	// Health
	router.Get("/health", HealthHandler)
	router.Get("/health/metrics", HealthHandler) // Prometheus style metrics?
}
