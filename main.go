package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/api"
	"gitlab.com/music-library/music-api/version"
)

// Measure time
// start := time.Now()
// diff = time.Now().Sub(start)

func main() {
	version.PrintTitle()

	// Initiate Fiber web-server
	app := fiber.New()

	// Setup the router
	api.MockRoutes(app)

	// Listen
	log.Fatal(app.Listen(listenAddr()))
}

func listenAddr() string {
	port := getPort()
	return fmt.Sprintf("localhost:%s", port)
}

func getPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "3001"
	}

	return port
}
