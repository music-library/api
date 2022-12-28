package api

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.com/music-library/music-api/global"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		description        string
		route              string
		method             string // HTTP method
		expectedStatusCode int    // Expected HTTP status code
		filesCount         uint64
		res                fiber.Map
	}{
		// Add test cases here
		{
			description:        "error if track index is empty",
			route:              "/health",
			method:             "GET",
			expectedStatusCode: 500,
			filesCount:         0,
			res: fiber.Map{
				"message": "track index is empty",
				"ok":      false,
			},
		},
		{
			description:        "returns status 200",
			route:              "/health",
			method:             "GET",
			expectedStatusCode: 200,
			filesCount:         1,
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
			filesCount:         45678,
			res: fiber.Map{
				"message": "ok",
				"ok":      true,
			},
		},
	}

	// Define Fiber app.
	app := fiber.New()

	// Create route with GET method for test
	app.Get("/health", HealthHandler)

	// Iterate through single test cases
	for _, test := range tests {
		global.Index.FilesCount = test.filesCount

		// Create a new http request with the route from the test case
		req := httptest.NewRequest(test.method, test.route, nil)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, 1)

		res, _ := sonic.Marshal(test.res)
		body, _ := io.ReadAll(resp.Body)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, test.description)
		assert.Equalf(t, test.method, resp.Request.Method, test.description)
		assert.Equalf(t, res, body, test.description)
	}
}
