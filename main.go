package main

import (
	"os"
	"sync"
	"time"

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
	// Create data directory
	if _, err := os.Stat(global.DATA_DIR); os.IsNotExist(err) {
		os.Mkdir(global.DATA_DIR, 0755)
	}

	// Create music directory
	if _, err := os.Stat(global.MUSIC_DIR); os.IsNotExist(err) {
		os.Mkdir(global.MUSIC_DIR, 0755)
	}

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

	// Async index population (to prevent blocking the server)
	go (func() {
		// Populate the index
		global.Index.Populate(global.MUSIC_DIR)

		// Read metadata from cache
		indexCache := global.Cache.ReadAndParseMetadata()

		start := time.Now()
		var await sync.WaitGroup

		// Populate metadata
		for _, indexTrack := range global.Index.Tracks {
			await.Add(1)

			go (func(indexTrack *indexer.IndexTrack) {
				defer await.Done()

				// Check if track metadata is cached
				cachedTrack, isCached := indexCache.Tracks[indexTrack.Id]

				if isCached {
					indexTrack.IdAlbum = cachedTrack.IdAlbum
					indexTrack.Metadata = cachedTrack.Metadata
				} else {
					global.Index.PopulateFileMetadata(indexTrack)
				}

				// Cover
				if !global.Cache.Exists(indexTrack.IdAlbum + "/cover.jpg") {
					trackCover, _ := indexer.GetTrackCover(indexTrack.Path)

					if trackCover != nil {
						// Save to global Cache
						global.Cache.Add(indexTrack.IdAlbum, "cover.jpg", trackCover)
						indexer.ResizeTrackCover(indexTrack.IdAlbum, "600")
					}
				}
			})(indexTrack)
		}

		await.Wait()
		log.Debug("main/metadata took ", time.Since(start))

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
