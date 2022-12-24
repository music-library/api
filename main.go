package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/api"
	"gitlab.com/music-library/music-api/constants"
	"gitlab.com/music-library/music-api/version"
)

// Measure time
// start := time.Now()
// diff = time.Now().Sub(start)

func init() {
	MakeLogger(constants.LOG_FILE)
}

func main() {
	version.PrintTitle()

	// Initiate Fiber web-server
	app := fiber.New()

	// Setup the router
	api.ApiRoutes(app)

	log.Debug("Listening on: " + ListenAddr())

	// Listen
	log.Fatal(app.Listen(ListenAddr()))
}
