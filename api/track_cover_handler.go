package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/indexer"
)

func TrackCoverHandler(c *fiber.Ctx) error {
	start := time.Now()

	index := indexer.Index{
		Files: make(map[string]*indexer.IndexFile, 1000),
	}

	index.Populate(global.MUSIC_DIR)

	trackId := strings.ToLower(c.Params("id"))
	track, ok := index.Files[trackId]

	if !ok {
		log.Error("http/track/" + trackId + " track does not exist")
		return c.Status(500).Send([]byte("{}"))
		// @TODO: Send default cover
	}

	trackCover, trackCoverMimetype := indexer.GetTrackCover(track.Path)

	fmt.Println(time.Since(start))

	// Send response as JPEG
	c.Response().Header.Add("Content-Type", trackCoverMimetype)
	return c.Send(trackCover)
}
