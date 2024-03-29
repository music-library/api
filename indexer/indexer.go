package indexer

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/hmerritt/go-ngram"
	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var MusicLibIndex = IndexMany{
	DefaultKey: "main",
	Indexes:    make(map[string]*Index),
	Socket:     NewSocket(),
}

var IndexNgram = ngram.NgramIndex{
	NgramMap:   make(map[string]map[int]*ngram.IndexValue),
	IndexesMap: make(map[int]*ngram.IndexValue),
	Ngram:      3,
}

type IndexMany struct {
	// Default (fallback) `Indexes` map key.
	DefaultKey string
	// Store multiple indexes. An index is the meta needed for an entire music library.
	Indexes map[string]*Index
	// Websocket state
	Socket *Socket
}

// Tracks as array
// Object with key as track id, value as arr index
type Index struct {
	Id        string                             `json:"id"`
	Name      string                             `json:"name"`
	Libraries []config.ConfigurationMusicLibrary `json:"libraries"`
	Tracks    []*IndexTrack                      `json:"tracks"`
	TracksKey map[string]int                     `json:"tracks_map"`
	Albums    map[string][]string                `json:"albums"` // albums[id_album] = []TracksKey
	Decades   []string                           `json:"decades"`
	Genres    []string                           `json:"genres"`
}

type IndexTrack struct {
	Id       string    `json:"id"`
	IdAlbum  string    `json:"id_album"`
	Path     string    `json:"-"`
	Metadata *Metadata `json:"metadata"`
	Stats    *Stat     `json:"stats"`
}

func GetNewIndex(name string) Index {
	return Index{
		Id:        slug.Make(name),
		Name:      cases.Title(language.AmericanEnglish).String(strings.ToLower(name)),
		Libraries: config.Config.MusicLibraries,
		Tracks:    make([]*IndexTrack, 0, 5000),
		TracksKey: make(map[string]int, 5000),
		Albums:    make(map[string][]string, 500),
	}
}

// Safely get track from index
func (index *Index) Get(trackId string) (*IndexTrack, bool) {
	trackIndex, ok := index.TracksKey[trackId]

	if !ok || trackIndex >= len(index.Tracks) {
		return nil, false
	}

	return index.Tracks[trackIndex], true
}

// Populate File index with audio `IndexTrack` objects
func (index *Index) Populate(path string) {
	start := time.Now()
	log.Debug("index/populate start " + path)

	err := filepath.WalkDir(path, func(itemPath string, entry os.DirEntry, err error) error {
		if err != nil {
			log.Error("index/populate file failed to be indexed ", err)
			return nil
		}

		if !entry.IsDir() && IsFileAudio(itemPath) {
			itemId := HashString(itemPath)

			index.TracksKey[itemId] = len(index.Tracks)
			index.Tracks = append(index.Tracks, &IndexTrack{
				Id:       itemId,
				Path:     itemPath,
				Metadata: GetEmptyMetadata(),
				Stats:    GetEmptyStat(),
			})
		}

		return nil
	})

	if err != nil {
		log.Fatal("index/populate failed ", err)
	}

	log.Debug(fmt.Sprintf("index/populate end. indexed %d items in %s %s", len(index.Tracks), time.Since(start), path))
}

// Populate IndexTrack with actual metadata
func (index *Index) PopulateFileMetadata(indexTrack *IndexTrack) *IndexTrack {
	trackIndex := index.TracksKey[indexTrack.Id]

	if index.Tracks[trackIndex].Metadata.Title == "(unknown)" {
		metadata := GetTrackMetadata(indexTrack.Path)
		index.Tracks[trackIndex].Metadata = metadata
	}

	index.Tracks[trackIndex].IdAlbum = HashString(index.Tracks[trackIndex].Metadata.Album + index.Tracks[trackIndex].Metadata.AlbumArtist)

	return index.Tracks[trackIndex]
}

// Crawl each file in the index
func (index *Index) Crawl(callback func(IndexTrack)) {
	for _, v := range index.Tracks {
		callback(*v)
	}
}

func GetTrackNgramString(indexTrack *IndexTrack) string {
	return strings.ToLower(fmt.Sprintf("%s %s %s %s", indexTrack.Metadata.Album, indexTrack.Metadata.AlbumArtist, indexTrack.Metadata.Artist, indexTrack.Metadata.Title))
}

// Check if file is audio file
func IsFileAudio(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".flac" || ext == ".mp3" || ext == ".ogg" || ext == ".wav"
}

// Hash string using SHA1, returns a string.
func HashString(str string) string {
	// Store string in buffer array
	buf := []byte(str)

	// Hash string, returns buffer array
	hash_byte := sha1.Sum(buf)

	// Convert buffer into hex string
	hash := hex.EncodeToString(hash_byte[:])

	// Return hash as string
	return hash
}
