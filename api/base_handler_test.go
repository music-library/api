package api

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.com/music-library/music-api/version"
)

func TestBaseHandler(t *testing.T) {
	tests := []struct {
		description        string
		route              string
		method             string // HTTP method
		expectedStatusCode int    // Expected HTTP status code
		res                fiber.Map
	}{
		// Add test cases here
		{
			description:        "returns status 200",
			route:              "/",
			method:             "GET",
			expectedStatusCode: 200,
			res: fiber.Map{
				"message": "Hello, World 👋!",
				"version": version.Version,
				"routes": []string{
					"/tracks",
					"/track/:id",
					"/track/:id/audio",
					"/track/:id/cover/:size?",
					"/albums",
					"/health",
					"/health/metrics",
				},
			},
		},
	}

	// Define Fiber app.
	app := fiber.New(fiber.Config{
		AppName:     "music-api",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Create route with GET method for test
	app.Get("/", BaseHandler)

	// Iterate through single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest(test.method, test.route, nil)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, -1)

		res, _ := json.Marshal(test.res)
		body, _ := io.ReadAll(resp.Body)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, test.description)
		assert.Equalf(t, test.method, resp.Request.Method, test.description)
		assert.Equalf(t, res, body, test.description)
	}
}
