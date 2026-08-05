package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/api"
	"gitlab.com/music-library/music-api/config"
	"gitlab.com/music-library/music-api/indexer"
	"gitlab.com/music-library/music-api/version"
)

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
		AppName: "music-api",
	})

	// Middleware
	app.Use(cors.New())
	app.Use("/ws", api.WebsocketUpgradeMiddleware)
	app.Use(recover.New()) // Prevent crashes due to panics

	if config.Config.LogLevel == "debug" {
		app.Use(logger.New())
	}

	// Setup the router
	api.ApiRoutes(app)
	api.WebsocketEventHanders(api.WsHub)
	go api.WsHub.Run() // Start websocket

	// Index all libraries on startup.
	// Setup CRON job to reindex libraries periodically.
	schedule, err := newReindexScheduler(config.Config.ReIndexEvery, indexer.IndexAllLibraries)
	if err != nil {
		log.WithError(err).Fatal("failed to configure library reindex schedule")
	}
	schedule.Start()

	// Listen
	log.Info("music-api server listening on " + ListenAddr())
	log.Fatal(app.Listen(ListenAddr()))
}

func newReindexScheduler(reindexEvery string, task func()) (gocron.Scheduler, error) {
	interval, err := time.ParseDuration(reindexEvery)
	if err != nil {
		return nil, fmt.Errorf("parse REINDEX_EVERY %q: %w", reindexEvery, err)
	}

	schedule, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	if _, err := schedule.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(task),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	); err != nil {
		_ = schedule.Shutdown()
		return nil, fmt.Errorf("create reindex job: %w", err)
	}

	return schedule, nil
}
