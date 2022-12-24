package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/constants"
	"gitlab.com/music-library/music-api/indexer"
)

func TrackCoverHandler(c *fiber.Ctx) error {
	start := time.Now()

	index := indexer.Index{
		Files: make(map[string]*indexer.IndexFile, 1000),
	}

	index.Populate(constants.MUSIC_DIR)

	trackId := strings.ToLower(c.Params("id"))
	track := index.Files[trackId]
	trackCover, trackCoverMimetype := indexer.GetTrackCover(track.Path)

	fmt.Println(time.Since(start))

	// Send response as JPEG
	c.Response().Header.Add("Content-Type", trackCoverMimetype)
	return c.Send(trackCover)
}
