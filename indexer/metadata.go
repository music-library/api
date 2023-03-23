package indexer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/dhowden/tag"
	log "github.com/sirupsen/logrus"
	useCache "gitlab.com/music-library/music-api/cache"
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

func initRawMetadataFromGoLib(filePath string) tag.Metadata {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Error("index/metadata/GoLib failed to open file " + filePath)
	}

	defer file.Close()

	meta, err := tag.ReadFrom(io.ReadSeeker(file))

	if err != nil {
		log.Error("index/metadata/GoLib failed to extract metadata from " + filePath)
	}

	return meta
}

// Extract metadata from file using Go lib.
//
// @Note: Doesn't support duration.
func getRawMetadataFromGoLib(baseMeta *Metadata, filePath string) *Metadata {
	meta := initRawMetadataFromGoLib(filePath)

	trackNo, _ := meta.Track()
	baseMeta.Track = trackNo
	baseMeta.Title = meta.Title()
	baseMeta.Artist = meta.Artist()
	baseMeta.AlbumArtist = meta.AlbumArtist()
	baseMeta.Album = meta.Album()
	baseMeta.Year = fmt.Sprint(meta.Year())
	baseMeta.Genre = meta.Genre()
	baseMeta.Composer = meta.Composer()

	// @Note - Duration is not supported by dhowden/tag
	// baseMeta.Duration = meta.Duration()

	return baseMeta
}

// Extract metadata from file using *external* CLI tool.
//
// @Note: Supports duration. Seems to be ~10x slower than Go lib.
func getRawMetadataFromMediaInfoCli(baseMeta *Metadata, filePath string) (*Metadata, error) {
	mediainfoJSON, err := exec.Command("mediainfo", "--Output=JSON", filePath).CombinedOutput()

	if err != nil {
		log.Error("index/metadata/CLI failed to extract mediainfo metadata: `" + filePath + "`")
		return baseMeta, err
	}

	mediainfo := struct {
		Media struct {
			Track []struct {
				Track       string `json:"Track_Position"`
				Title       string `json:"Title"`
				Artist      string `json:"Performer"`
				AlbumArtist string `json:"Album_Performer"`
				Album       string `json:"Album"`
				Year        string `json:"Recorded_Date"`
				Genre       string `json:"Genre"`
				Composer    string `json:"Composer"`
				Duration    string `json:"Duration"`
			} `json:"track"`
		} `json:"media"`
	}{}

	// Parse JSON
	err = sonic.Unmarshal(mediainfoJSON, &mediainfo)

	if err != nil || len(mediainfo.Media.Track) == 0 {
		log.Error("index/metadata/CLI failed to unmarshal mediainfo json: `" + filePath + "`")
		return baseMeta, err
	}

	track := mediainfo.Media.Track[0]
	durationFloat, _ := strconv.ParseFloat(track.Duration, 32)
	trackNo, _ := strconv.Atoi(strings.TrimLeft(track.Track, "0")) // Trim leading zeros
	baseMeta.Track = trackNo
	baseMeta.Title = track.Title
	baseMeta.Artist = track.Artist
	baseMeta.AlbumArtist = track.AlbumArtist
	baseMeta.Album = track.Album
	baseMeta.Year = track.Year
	baseMeta.Genre = track.Genre
	baseMeta.Composer = track.Composer
	baseMeta.Duration = int(durationFloat)

	return baseMeta, nil
}

func refineRawMetadata(meta *Metadata, filePath string) *Metadata {
	if len(meta.Title) == 0 {
		// Use filename as title if no title is found
		meta.Title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

		// Attempt to clean up title (worth a shot)
		meta.Title = strings.ReplaceAll(meta.Title, "_", " ")
		meta.Title = strings.ReplaceAll(meta.Title, "-", " ")
		meta.Title = strings.ReplaceAll(meta.Title, ".", " ")
		meta.Title = strings.ReplaceAll(meta.Title, "   ", " ")
		meta.Title = strings.ReplaceAll(meta.Title, "  ", " ")
		meta.Title = cases.Title(language.AmericanEnglish).String(strings.ToLower(meta.Title))
	}

	if len(meta.Artist) == 0 {
		meta.Artist = "~"
	}

	if len(meta.AlbumArtist) == 0 {
		meta.AlbumArtist = "~"
	}

	if len(meta.Album) == 0 {
		meta.Album = "~"
	}

	if len(meta.Year) < 4 || len(meta.Year) > 4 {
		meta.Year = "~"
	}

	if len(meta.Year) == 4 {
		meta.Decade = meta.Year[:3] + "0"
	}

	if len(meta.Genre) == 0 {
		meta.Genre = "~"
	} else {
		meta.Genre = cases.Title(language.AmericanEnglish).String(strings.ToLower(meta.Genre))
	}

	if len(meta.Composer) == 0 {
		meta.Composer = "~"
	}

	return meta
}

func GetTrackMetadata(filePath string) *Metadata {
	meta := GetEmptyMetadata()
	err := error(nil)

	meta, err = getRawMetadataFromMediaInfoCli(meta, filePath)

	if err != nil {
		meta = getRawMetadataFromGoLib(meta, filePath)
	}

	meta = refineRawMetadata(meta, filePath)

	return meta
}

// Returns cover image as byte array, and mime type
func GetTrackCover(filePath string) ([]byte, string) {
	meta := initRawMetadataFromGoLib(filePath)

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

func ResizeTrackCover(idAlbum string, size string, cache useCache.Cache) (string, error) {
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

func ReadAndParseMetadata(cache useCache.Cache) *Index {
	metadataRaw, err := cache.Read("metadata.json")
	indexCache := &Index{}

	if err != nil {
		log.Error("cache/parse/metadata failed to read cache file metadata.json")
	}

	if err == nil {
		sonic.Unmarshal(metadataRaw, indexCache)
	}

	return indexCache
}
