package api

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.com/music-library/music-api/indexer"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		description        string
		route              string
		method             string // HTTP method
		expectedStatusCode int    // Expected HTTP status code
		tracksCount        uint64
		indexNames         []string
		res                fiber.Map
	}{
		// Add test cases here
		{
			description:        "error if track index is empty",
			route:              "/health",
			method:             "GET",
			expectedStatusCode: 500,
			tracksCount:        0,
			indexNames:         []string{"empty", "empty2", "empty3"},
			res: fiber.Map{
				"message": "empty track index is empty",
				"ok":      false,
			},
		},
		{
			description:        "returns status 200",
			route:              "/health",
			method:             "GET",
			expectedStatusCode: 200,
			tracksCount:        1,
			indexNames:         []string{"main", "second", "third"},
			res: fiber.Map{
				"message": "ok",
				"ok":      true,
			},
		},
		{
			description:        "returns status 200",
			route:              "/health",
			method:             "GET",
			expectedStatusCode: 200,
			tracksCount:        456,
			indexNames:         []string{"main", "456", "rand456"},
			res: fiber.Map{
				"message": "ok",
				"ok":      true,
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
	app.Get("/health", HealthHandler)

	// Iterate through single test cases
	for _, test := range tests {
		indexMany := *indexer.TestGenerateIndexMany(test.indexNames, test.tracksCount)
		indexer.MusicLibIndex = indexMany

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
