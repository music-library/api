package main

import (
	"os"
	"sort"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hmerritt/go-ngram"
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
	//
	// Uses custom JSON encoding as recommended: https://docs.gofiber.io/guide/faster-fiber
	app := fiber.New(fiber.Config{
		AppName:     "music-api",
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Middleware
	app.Use(recover.New()) // Prevent crashes due to panics

	if config.Config.LogLevel == "debug" {
		app.Use(logger.New())
	}

	// Setup the router
	api.ApiRoutes(app)

	// Async index population (to prevent blocking the server)
	go (func() {
		// Populate the index
		global.Index.Populate(config.Config.MusicDir)

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
				cachedTrackIndex, isCached := indexCache.TracksKey[indexTrack.Id]

				if isCached {
					cachedTrack := indexCache.Tracks[cachedTrackIndex]
					indexTrack.IdAlbum = cachedTrack.IdAlbum
					indexTrack.Metadata = cachedTrack.Metadata
					indexTrack.Stats = cachedTrack.Stats
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

		// Second sync pass
		decadeKeys := make(map[string]bool)
		genresKeys := make(map[string]bool)

		for index, track := range global.Index.Tracks {
			// ngram index
			global.IndexNgram.Add(indexer.GetTrackNgramString(track), ngram.NewIndexValue(index, track))

			// albums
			_, ok := global.Index.Albums[track.IdAlbum]
			if !ok {
				global.Index.Albums[track.IdAlbum] = make([]string, 0, 20)
			}
			global.Index.Albums[track.IdAlbum] = append(global.Index.Albums[track.IdAlbum], track.Id)

			// decades
			if _, ok := decadeKeys[track.Metadata.Decade]; !ok {
				decade := track.Metadata.Decade
				decadeKeys[decade] = true

				if len(decade) == 4 {
					global.Index.Decades = append(global.Index.Decades, decade)
				}
			}

			// genres
			if _, ok := genresKeys[track.Metadata.Genre]; !ok {
				genre := track.Metadata.Genre
				genresKeys[genre] = true

				if len(genre) > 0 {
					global.Index.Genres = append(global.Index.Genres, genre)
				}
			}
		}

		sort.Slice(global.Index.Decades, func(i, j int) bool {
			return global.Index.Decades[i] < global.Index.Decades[j]
		})

		sort.Slice(global.Index.Genres, func(i, j int) bool {
			return global.Index.Genres[i] < global.Index.Genres[j]
		})

		log.Info("main/metadata took ", time.Since(start))

		// Cache metadata
		metadataJSON, err := sonic.Marshal(global.Index)

		if err != nil {
			log.Error("main/metadata/cache failed to marshal metadata ", err)
		}

		global.Cache.Replace(".", "metadata.json", metadataJSON)
	})()

	// Listen
	log.Info("music-api server listening on " + ListenAddr())
	log.Fatal(app.Listen(ListenAddr()))
}
