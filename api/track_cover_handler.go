package api

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/global"
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
		// Check if cover exists in global Cache
		imgSizePath := fmt.Sprintf("%s/%s.jpg", track.IdAlbum, c.Params("size"))

		if !global.Cache.Exists(imgSizePath) {
			exec.Command("vipsthumbnail", global.Cache.FilePath(imgPath), "--size", c.Params("size"), "-o", fmt.Sprintf("%s[Q=90]", global.Cache.FilePath(imgSizePath))).Run()
		}

		return c.SendFile(global.Cache.FilePath(imgSizePath))
	}

	return c.SendFile(global.Cache.FilePath(imgPath))
}
