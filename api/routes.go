package api

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(router fiber.Router) {
	// Track
	router.Get("/tracks", MockHandler)
	router.Get("/track/:id", MockHandler)
	router.Get("/track/:id/audio", MockHandler)
	router.Get("/track/:id/cover/:size?", MockHandler)

	// Album
	router.Get("/albums", MockHandler)

	// Health
	router.Get("/health", MockHandler)
	router.Get("/health/metrics", MockHandler) // Prometheus style metrics?
}
