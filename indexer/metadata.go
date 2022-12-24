package indexer

import (
	"io"
	"os"

	"github.com/dhowden/tag"
	log "github.com/sirupsen/logrus"
)

type Metadata struct {
	Track        int
	Title        string
	Artist       string
	Album_artist string
	Album        string
	Year         int
	Genre        string
	Composer     string
	Duration     int // in seconds
	// Cover     not here -> stored in cache
}

func GetEmptyMetadata() *Metadata {
	return &Metadata{
		Track:        1,
		Title:        "(unknown)",
		Artist:       "~",
		Album_artist: "~",
		Album:        "~",
		Year:         0,
		Genre:        "~",
		Composer:     "~",
		Duration:     0,
	}
}

func getRawMetadata(filePath string) tag.Metadata {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Error("index/metadata failed to open file " + filePath)
	}

	meta, err := tag.ReadFrom(io.ReadSeeker(file))

	if err != nil {
		log.Error("index/metadata failed to extract metadata from " + filePath)
	}

	return meta
}

func GetTrackMetadata(filePath string) *Metadata {
	meta := getRawMetadata(filePath)

	track, _ := meta.Track()

	return &Metadata{
		Track:        track,
		Title:        meta.Title(),
		Artist:       meta.Artist(),
		Album_artist: meta.AlbumArtist(),
		Album:        meta.Album(),
		Year:         meta.Year(),
		Genre:        meta.Genre(),
		Composer:     meta.Composer(),
		Duration:     0,
		// Duration:     meta.Duration(),
	}
}

// Returns cover image as byte array, and mime type
func GetTrackCover(filePath string) ([]byte, string) {
	meta := getRawMetadata(filePath)

	picture := meta.Picture()

	if picture == nil {
		log.Debug("index/metadata track has no cover image " + filePath)
		return nil, "image/jpeg"
	}

	mimeType := picture.MIMEType

	if len(mimeType) == 0 {
		mimeType = "image/jpeg"
	}

	return picture.Data, mimeType
}
