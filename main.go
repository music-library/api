package main

import (
	"os"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/api"
	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/indexer"
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
		AppName:     "music-api",
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Middleware
	app.Use(recover.New()) // Prevent crashes due to panics

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

		var await sync.WaitGroup

		// Populate metadata
		for _, indexFile := range global.Index.Files {
			await.Add(1)

			go (func(indexFile *indexer.IndexFile) {
				defer await.Done()
				global.Index.PopulateFileMetadata(indexFile)

				// Cover
				if !global.Cache.Exists(indexFile.IdAlbum + "/cover.jpg") {
					trackCover, _ := indexer.GetTrackCover(indexFile.Path)

					if trackCover != nil {
						// Save to global Cache
						global.Cache.Add(indexFile.IdAlbum, "cover.jpg", trackCover)
					}
				}
			})(indexFile)
		}

		await.Wait()

		// Cache metadata
		metadataJSON, err := sonic.Marshal(global.Index)

		if err != nil {
			log.Error("main/metadata/cache failed to marshal metadata ", err)
		}

		global.Cache.Replace(".", "metadata.json", metadataJSON)
	})()

	// Listen
	log.Debug("music-api server listening on " + ListenAddr())
	log.Fatal(app.Listen(ListenAddr()))
}
