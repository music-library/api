package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/indexer"
	"gitlab.com/music-library/music-api/static"
)

func TrackCoverHandler(c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "image/jpg")

	trackId := strings.ToLower(c.Params("id"))
	track, ok := global.Index.Get(trackId)
	cache := indexer.GetCache()

	if !ok {
		log.Error("http/track/" + trackId + "/cover track does not exist")
		return c.Send(getDefaultCover())
	}

	imgPath := fmt.Sprintf("%s/cover.jpg", track.IdAlbum)

	if !cache.Exists(imgPath) {
		log.Error("http/track/" + trackId + "/cover cover does not exist")
		return c.Send(getDefaultCover())
	}

	if len(c.Params("size")) > 0 {
		imgResizePath, _ := indexer.ResizeTrackCover(track.IdAlbum, c.Params("size"))
		return c.SendFile(cache.FilePath(imgResizePath))
	}

	return c.SendFile(cache.FilePath(imgPath))
}

func getDefaultCover() []byte {
	data, err := static.Images.ReadFile("images/image-placeholder.jpg")

	if err != nil {
		log.Error("http/track/cover failed to getDefaultCover() " + err.Error())
	}

	return data
}
