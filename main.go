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

	// Create data directory
	if _, err := os.Stat(constants.DATA_DIR); os.IsNotExist(err) {
		log.Debug("creating data directory " + constants.DATA_DIR)
		os.Mkdir(constants.DATA_DIR, 0755)
	}

	// Listen
	log.Debug("music-api server listening on " + ListenAddr())
	log.Fatal(app.Listen(ListenAddr()))
}
