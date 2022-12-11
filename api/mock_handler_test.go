package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestMockHandler(t *testing.T) {
	tests := []struct {
		description        string
		route              string
		method             string // HTTP method
		bodyLength         int
		expectedStatusCode int // Expected HTTP status code
	}{
		// Add test cases here
		{
			description:        "get HTTP status 200",
			route:              "/",
			method:             "GET",
			bodyLength:         0,
			expectedStatusCode: 200,
		},
		{
			description:        "get HTTP status 200",
			route:              "/anything",
			method:             "POST",
			bodyLength:         0,
			expectedStatusCode: 200,
		},
		{
			description:        "get HTTP status 200",
			route:              "/anything",
			method:             "DELETE",
			bodyLength:         0,
			expectedStatusCode: 200,
		},
	}

	// Define Fiber app.
	app := fiber.New()

	// Create route with GET method for test
	app.Get("/*", MockHandler)
	app.Post("/*", MockHandler)
	app.Delete("/*", MockHandler)

	// Iterate through single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest(test.method, test.route, nil)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, 1)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, test.description)
		assert.Equalf(t, test.method, resp.Request.Method, test.description)
	}
}
