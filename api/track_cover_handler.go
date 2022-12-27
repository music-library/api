package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/indexer"
)

func TrackCoverHandler(c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "image/jpg")

	trackId := strings.ToLower(c.Params("id"))
	track, ok := global.Index.Files[trackId]

	if !ok {
		log.Error("http/track/" + trackId + "/cover track does not exist")
		return Error(c, 500, "track does not exist")
		// @TODO: Send default cover
	}

	imgPath := fmt.Sprintf("%s/cover.jpg", track.IdAlbum)

	if !global.Cache.Exists(imgPath) {
		log.Error("http/track/" + trackId + "/cover cover does not exist")
		return Error(c, 500, "cover does not exist")
		// @TODO: Send default cover
	}

	if len(c.Params("size")) > 0 {
		imgResizePath, _ := indexer.ResizeTrackCover(track.IdAlbum, c.Params("size"))
		return c.SendFile(global.Cache.FilePath(imgResizePath))
	}

	return c.SendFile(global.Cache.FilePath(imgPath))
}
