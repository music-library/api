package api

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/indexer"
)

func TrackAudioHandler(c *fiber.Ctx) error {
	libId := c.Locals("libId").(string)
	trackId := strings.ToLower(c.Params("id"))
	track, ok := indexer.MusicLibIndex.Indexes[libId].Get(trackId)

	if !ok {
		log.Error("http/track/" + trackId + "/audio track does not exist")
		return Error(c, 404, "track does not exist")
	}

	// Open track file for reading
	file, fileStat, err := openTrack(trackId, track.Path)
	if err != nil {
		return Error(c, 500, "track failed to play")
	}
	defer file.Close()

	totalSize := fileStat.Size()
	lastModified := fileStat.ModTime()

	mimeType := mime.TypeByExtension(filepath.Ext(track.Path))

	if mimeType == "" || mimeType == "audio/x-flac" {
		mimeType = "audio/flac"
	}

	c.Set(fiber.HeaderContentType, mimeType)
	c.Set(fiber.HeaderContentLength, fmt.Sprint(totalSize))
	c.Set(fiber.HeaderLastModified, lastModified.Format(http.TimeFormat))
	c.Set(fiber.HeaderCacheControl, "public, max-age=31536000")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("inline; filename=%s - %s%s", track.Metadata.Artist, track.Metadata.Title, filepath.Ext(track.Path)))

	// Accept range requests
	reqRange := c.Get(fiber.HeaderRange)

	// Only increment times played at the initial request (not range requests)
	if reqRange == "" || reqRange == "bytes=0-" {
		// Update track stats
		track.Stats.TimesPlayed += 1
		track.Stats.LastPlayed = time.Now().Unix()
	}

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
		buffer := make([]byte, chunksize)
		bytesread, err := file.ReadAt(buffer, start)

		if err != nil {
			log.Error("http/track/" + trackId + "/audio track file failed to read correctly")
		}

		c.Set(fiber.HeaderAcceptRanges, "bytes")
		c.Set(fiber.HeaderContentLength, fmt.Sprint(chunksize))
		c.Set(fiber.HeaderContentRange, fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
		c.Status(fiber.StatusPartialContent)

		return c.Send(buffer[:bytesread])
	}

	return c.SendFile(track.Path)
}

// Open track and get file size
func openTrack(trackId, trackPath string) (*os.File, fs.FileInfo, error) {
	file, err := os.Open(trackPath)

	if err != nil {
		log.Error("http/track/" + trackId + "/audio track file failed to open")
		return nil, nil, err
	}

	fileStat, err := file.Stat()

	if err != nil {
		log.Error("http/track/" + trackId + "/audio track file failed to get file.stat")
		return nil, nil, err
	}

	return file, fileStat, nil
}
