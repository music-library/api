package api

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strconv"
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

	totalSize := trackFileInfo.Size()
	c.Set(fiber.HeaderContentType, mimeType)
	c.Set(fiber.HeaderContentLength, fmt.Sprint(totalSize))

	// Accept range requests
	reqRange := c.Get(fiber.HeaderRange)
	if reqRange != "" {
		parts := strings.Split(strings.Replace(reqRange, "bytes=", "", 1), "-")
		partialstart := parts[0]
		partialend := parts[1]

		start, err := strconv.ParseInt(partialstart, 10, 64)
		if err != nil {
			start = 0
		}

		end, err := strconv.ParseInt(partialend, 10, 64)
		if err != nil || partialend != "" {
			end = totalSize - 1
		}

		chunksize := end - start + 1

		file, err := os.Open(track.Path)
		if err != nil {
			log.Error("http/track/" + trackId + "/audio track failed to play")
		}

		defer file.Close()
		file.Seek(start, 0)
		fileReader := io.LimitReader(file, totalSize-start)

		c.Set(fiber.HeaderAcceptRanges, "bytes")
		c.Set(fiber.HeaderContentRange, fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
		c.Set(fiber.HeaderContentLength, fmt.Sprint(chunksize))

		return c.SendStream(fileReader)
	}

	return c.SendFile(track.Path, false)
}
