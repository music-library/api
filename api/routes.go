package api

import (
	"github.com/gofiber/fiber/v2"
)

//
// MOCK DUMMY API
//
// Returns JSON of request
func MockRoutes(router fiber.Router) {
	// Listen on for any route across HTTP methods

	// Listen on all HTTP methods, for any route
	router.All("/*", MockHandler)
}
