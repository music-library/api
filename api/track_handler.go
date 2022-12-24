package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/constants"
	"gitlab.com/music-library/music-api/indexer"
)

func TrackHandler(c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "application/json")

	start := time.Now()

	index := indexer.Index{
		Files: make(map[string]*indexer.IndexFile, 1000),
	}

	index.Populate(constants.MUSIC_DIR)

	trackId := strings.ToLower(c.Params("id"))
	track, ok := index.Files[trackId]

	if !ok {
		log.Error("http/track/" + trackId + " track does not exist")
		return c.Status(500).Send([]byte("{}"))
	}

	index.PopulateFileMetadata(track)

	trackJSON, err := json.Marshal(track)

	if err != nil {
		log.Error("http/track/" + trackId + " failed to marshal track json")
		return c.Status(500).Send([]byte("{}"))
	}

	fmt.Println(time.Since(start))

	return c.Send(trackJSON)
}
