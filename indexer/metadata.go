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
	Duration     string
	Raw          map[string]interface{}
	// Cover     not here -> stored in cache
}

func (index *Index) GetMetadata(filePath string) *Metadata {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Error("index/metadata failed to open file " + filePath)
	}

	meta, err := tag.ReadFrom(io.ReadSeeker(file))

	if err != nil {
		log.Error("index/metadata failed to extract metadata from " + filePath)
	}

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
		Raw:          meta.Raw(),
		// Duration:     meta.Duration(),
	}
}
