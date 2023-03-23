package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/indexer"
)

func TrackHandler(c *fiber.Ctx) error {
	libId := c.Locals("libId").(string)
	trackId := strings.ToLower(c.Params("id"))
	track, ok := indexer.MusicLibIndex.Indexes[libId].Get(trackId)

	if !ok {
		log.Error("http/track/" + trackId + " track does not exist")
		return Error(c, 404, "track does not exist")
	}

	return c.JSON(track)
}
