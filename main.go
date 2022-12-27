package main

import (
	"os"

	"github.com/bytedance/sonic"
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	//
	// Uses custom JSON encoding as recommended: https://docs.gofiber.io/guide/faster-fiber
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Middleware
	if global.LOG_LEVEL == "debug" {
		app.Use(logger.New())
	}

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

		// Populate metadata
		for _, indexFile := range global.Index.Files {
			go (func() {
				global.Index.PopulateFileMetadata(indexFile)
			})()
		}
	})()

	// Listen
	log.Debug("music-api server listening on " + ListenAddr())
	log.Fatal(app.Listen(ListenAddr()))
}
