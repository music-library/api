package api

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(router fiber.Router) {
	router.All("/", BaseHandler)

	// Track
	router.Get("/tracks", TracksHandler)
	router.Get("/tracks/search/:query", SearchHandler)

	router.Get("/track/:id", TrackHandler)
	router.Get("/track/:id/audio", TrackAudioHandler)
	router.Get("/track/:id/cover/:size?", TrackCoverHandler)

	// LEGACY
	router.Get("/tracks/:id", TrackHandler)
	router.Get("/tracks/:id/audio", TrackAudioHandler)
	router.Get("/tracks/:id/cover/:size?", TrackCoverHandler)

	// Album
	router.Get("/albums", AlbumsHandler)

	// Health
	router.Get("/health", HealthHandler)
	router.Get("/health/metrics", HealthHandler) // Prometheus style metrics?
}
