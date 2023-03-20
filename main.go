package main

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/api"
	"gitlab.com/music-library/music-api/config"
	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/indexer"
	"gitlab.com/music-library/music-api/version"
)

// Measure time
// start := time.Now()
// time.Since(start)
// diff = time.Now().Sub(start)

func init() {
	// Create data directory
	if _, err := os.Stat(config.Config.DataDir); os.IsNotExist(err) {
		os.Mkdir(config.Config.DataDir, 0755)
	}

	// Create music directory
	if _, err := os.Stat(config.Config.MusicDir); os.IsNotExist(err) {
		os.Mkdir(config.Config.MusicDir, 0755)
	}

	MakeLogger(config.Config.LogFile)
}

func main() {
	version.PrintTitle()

	// Initiate Fiber web-server
	// Uses custom JSON encoding as recommended: https://docs.gofiber.io/guide/faster-fiber
	app := fiber.New(fiber.Config{
		AppName:     "music-api",
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Middleware
	app.Use(cors.New())
	app.Use(recover.New()) // Prevent crashes due to panics

	if config.Config.LogLevel == "debug" {
		app.Use(logger.New())
	}

	// Setup the router
	api.ApiRoutes(app)

	// Async index population (to prevent blocking the server)
	go (func() {
		// Index all music libraries
		for _, musicLibConfig := range config.Config.MusicLibraries {
			mainIndex := indexer.BootstrapIndex(musicLibConfig.Name, musicLibConfig.Path)
			global.IndexMany.Indexes[mainIndex.Id] = mainIndex
		}
	})()

	// Listen
	log.Info("music-api server listening on " + ListenAddr())
	log.Fatal(app.Listen(ListenAddr()))
}
