package indexer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Stat struct {
	TimesPlayed uint64 `json:"timesPlayed"`
	LastPlayed  int64  `json:"lastPlayed"`
}

func GetEmptyStat() *Stat {
	return &Stat{
		TimesPlayed: 0,
		LastPlayed:  -1,
	}
}

type Metadata struct {
	Track       int    `json:"track"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	AlbumArtist string `json:"album_artist"`
	Album       string `json:"album"`
	Year        string `json:"year"`
	Decade      string `json:"decade"`
	Genre       string `json:"genre"`
	Composer    string `json:"composer"`
	Duration    int    `json:"duration"` // in seconds
	// Cover     not here -> stored in cache
}

func GetEmptyMetadata() *Metadata {
	return &Metadata{
		Track:       0,
		Title:       "(unknown)",
		Artist:      "~",
		AlbumArtist: "~",
		Album:       "~",
		Year:        "~",
		Decade:      "~",
		Genre:       "~",
		Composer:    "~",
		Duration:    0,
	}
}

func getRawMetadata(filePath string) tag.Metadata {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Error("index/metadata failed to open file " + filePath)
	}

	defer file.Close()

	meta, err := tag.ReadFrom(io.ReadSeeker(file))

	if err != nil {
		log.Error("index/metadata failed to extract metadata from " + filePath)
	}

	return meta
}

func GetTrackMetadata(filePath string) *Metadata {
	baseMeta := GetEmptyMetadata()
	meta := getRawMetadata(filePath)

	track, _ := meta.Track()
	baseMeta.Track = track

	title := meta.Title()
	if len(title) == 0 {
		// Use filename as title if no title is found
		title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

		// Attempt to clean up title (worth a shot)
		title = strings.ReplaceAll(title, "_", " ")
		title = strings.ReplaceAll(title, "-", " ")
		title = strings.ReplaceAll(title, ".", " ")
		title = strings.ReplaceAll(title, "   ", " ")
		title = strings.ReplaceAll(title, "  ", " ")
		title = cases.Title(language.AmericanEnglish).String(strings.ToLower(title))
	}
	baseMeta.Title = title

	artist := meta.Artist()
	if len(artist) == 0 {
		artist = "~"
	}
	baseMeta.Artist = artist

	albumArtist := meta.AlbumArtist()
	if len(albumArtist) == 0 {
		albumArtist = "~"
	}
	baseMeta.AlbumArtist = albumArtist

	album := meta.Album()
	if len(album) == 0 {
		album = "~"
	}
	baseMeta.Album = album

	year := fmt.Sprint(meta.Year())
	if len(year) < 4 {
		year = "~"
	}
	baseMeta.Year = year

	if len(year) == 4 {
		baseMeta.Decade = year[:3] + "0"
	}

	genre := cases.Title(language.AmericanEnglish).String(strings.ToLower(meta.Genre()))
	if len(genre) == 0 {
		genre = "~"
	}
	baseMeta.Genre = genre

	composer := meta.Composer()
	if len(composer) == 0 {
		composer = "~"
	}
	baseMeta.Composer = composer

	// @TODO
	// baseMeta.Duration = meta.Duration()

	return baseMeta
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

func ResizeTrackCover(idAlbum string, size string) (string, error) {
	cache := Cache{
		Path: "./data",
	}
	imgPath := fmt.Sprintf("%s/cover.jpg", idAlbum)
	imgResizePath := fmt.Sprintf("%s/%s.jpg", idAlbum, size)

	if !cache.Exists(imgResizePath) {
		err := exec.Command("vipsthumbnail", cache.FilePath(imgPath), "--size", size, "-o", fmt.Sprintf("%s[Q=90]", cache.FilePath(imgResizePath))).Run()

		if err != nil {
			log.Error("index/metadata failed to resize album cover image " + idAlbum + " err: " + err.Error())
			return imgPath, err
		}
	}

	return imgResizePath, nil
}
