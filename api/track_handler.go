package api

import (
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"gitlab.com/music-library/music-api/global"
)

func TrackHandler(c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "application/json")

	trackId := strings.ToLower(c.Params("id"))
	track, ok := global.Index.Files[trackId]

	if !ok {
		log.Error("http/track/" + trackId + " track does not exist")
		return c.Status(500).Send([]byte("{}"))
	}

	global.Index.PopulateFileMetadata(track)

	trackJSON, err := sonic.Marshal(track)

	if err != nil {
		log.Error("http/track/" + trackId + " failed to marshal track json")
		return c.Status(500).Send([]byte("{}"))
	}

	return c.Send(trackJSON)
}
