package api

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
	trackMetaFriendlyName := regexp.MustCompile(`[^a-zA-Z0-9\s\-\_\.]`).ReplaceAllString(fmt.Sprintf("%s - %s%s", track.Metadata.Artist, track.Metadata.Title, filepath.Ext(track.Path)), "")

	mimeType := mime.TypeByExtension(filepath.Ext(track.Path))

	if mimeType == "" || mimeType == "audio/x-flac" {
		mimeType = "audio/flac"
	}

	c.Set(fiber.HeaderContentType, mimeType)
	c.Set(fiber.HeaderContentLength, fmt.Sprint(totalSize))
	c.Set(fiber.HeaderLastModified, lastModified.Format(http.TimeFormat))
	c.Set(fiber.HeaderCacheControl, "public, max-age=31536000")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("inline; filename=%s", trackMetaFriendlyName))

	// Accept range requests
	reqRange := c.Get(fiber.HeaderRange)

	// Only increment times played at the initial request (not range requests)
	if reqRange == "" || reqRange == "bytes=0-" {
		// Update track stats
		track.Stats.TimesPlayed += 1
		track.Stats.LastPlayed = time.Now().Unix()
	}

	if reqRange == "" {
		// No Range header, send the entire file
		_, err = io.Copy(c, file)
		if err != nil {
			log.Error("http/track/" + trackId + "/audio track failed to copy")
		}
		return err
	}

	// Parse the Range header
	rangeValues := strings.Split(strings.Split(reqRange, "=")[1], "-")
	startByte, _ := strconv.ParseInt(rangeValues[0], 10, 64)
	endByte := totalSize - 1
	chunksize := endByte - startByte + 1

	if len(rangeValues) == 2 && rangeValues[1] != "" {
		endByte, _ = strconv.ParseInt(rangeValues[1], 10, 64)
		chunksize = endByte - startByte + 1
	}

	c.Status(fiber.StatusPartialContent)
	c.Set(fiber.HeaderAcceptRanges, "bytes")
	c.Set(fiber.HeaderContentLength, strconv.FormatInt(chunksize, 10))
	c.Set(fiber.HeaderContentRange, fmt.Sprintf("bytes %d-%d/%d", startByte, endByte, totalSize))

	// Seek to the start byte and stream the audio
	_, err = file.Seek(startByte, io.SeekStart)
	if err != nil {
		log.Error("http/track/" + trackId + "/audio track failed seek to the start byte")
		return Error(c, 500, "track failed to seek to the start byte")
	}

	_, err = io.Copy(c, io.LimitReader(file, chunksize))
	if err != nil {
		log.Error("http/track/" + trackId + "/audio track failed to copy")
	}
	return err
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
