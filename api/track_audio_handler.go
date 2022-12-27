package api

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"gitlab.com/music-library/music-api/global"
)

func TrackAudioHandler(c *fiber.Ctx) error {
	trackId := strings.ToLower(c.Params("id"))
	track, ok := global.Index.Files[trackId]

	if !ok {
		log.Error("http/track/" + trackId + "/audio track does not exist")
		return Error(c, 404, "track does not exist")
	}

	trackFileInfo, err := os.Stat(track.Path)

	if err != nil {
		log.Error("http/track/" + trackId + "/audio track failed to play")
		return Error(c, 500, "track failed to play")
	}

	mimeType := mime.TypeByExtension(filepath.Ext(track.Path))

	if mimeType == "" || mimeType == "audio/x-flac" {
		mimeType = "audio/flac"
	}

	c.Set(fiber.HeaderContentLength, fmt.Sprint(trackFileInfo.Size()))
	c.Set(fiber.HeaderContentType, mimeType)
	return c.SendFile(track.Path, false)
}
