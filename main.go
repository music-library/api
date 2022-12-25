package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/api"
	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/version"
)

// Measure time
// start := time.Now()
// time.Since(start)
// diff = time.Now().Sub(start)

func init() {
	MakeLogger(global.LOG_FILE)
}

func main() {
	version.PrintTitle()

	// Initiate Fiber web-server
	app := fiber.New()

	// Setup the router
	api.ApiRoutes(app)

	// Create data directory
	if _, err := os.Stat(global.DATA_DIR); os.IsNotExist(err) {
		log.Debug("creating data directory " + global.DATA_DIR)
		os.Mkdir(global.DATA_DIR, 0755)
	}

	// Populate the index
	go (func() {
		global.Index.Populate(global.MUSIC_DIR)
	})()

	// Listen
	log.Debug("music-api server listening on " + ListenAddr())
	log.Fatal(app.Listen(ListenAddr()))
}
