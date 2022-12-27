package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/h2non/bimg"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/indexer"
)

func TrackCoverHandler(c *fiber.Ctx) error {
	trackId := strings.ToLower(c.Params("id"))
	track, ok := global.Index.Files[trackId]

	if !ok {
		log.Error("http/track/" + trackId + "/cover track does not exist")
		return Error(c, 500, "track does not exist")
		// @TODO: Send default cover
	}

	trackCover, trackCoverMimetype := indexer.GetTrackCover(track.Path)

	newImage, err := bimg.NewImage(trackCover).Process(bimg.Options{
		Width:   200,
		Height:  200,
		Quality: 90,
		Crop:    false,
	})

	if err != nil {
		log.Error("http/track/" + trackId + "/cover failed resizing cover")
		return Error(c, 500, "failed resizing cover")
		// @TODO: Send default cover
	}

	// Send response as JPEG
	c.Response().Header.Add("Content-Type", trackCoverMimetype)
	return c.Send(newImage)
}
