package indexer

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type IndexTrack struct {
	Id       string    `json:"id"`
	IdAlbum  string    `json:"id_album"`
	Path     string    `json:"path"`
	Metadata *Metadata `json:"metadata"`
	// Stats    *Stat     `json:"stats"`
}

type Index struct {
	Tracks      map[string]*IndexTrack
	TracksCount uint64
	// Albums [][]string // Slice of IndexTrack.Id
	// Genres []string
	// Decades []string
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

			index.TracksCount += 1
			index.Tracks[itemId] = &IndexTrack{
				Id:       itemId,
				Path:     itemPath,
				Metadata: GetEmptyMetadata(),
				// Metadata: GetTrackMetadata(itemPath),
			}
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
	if index.Tracks[indexTrack.Id].Metadata.Title == "(unknown)" {
		metadata := GetTrackMetadata(indexTrack.Path)
		index.Tracks[indexTrack.Id].Metadata = metadata
	}

	index.Tracks[indexTrack.Id].IdAlbum = HashString(index.Tracks[indexTrack.Id].Metadata.Album + index.Tracks[indexTrack.Id].Metadata.Album_artist)

	return index.Tracks[indexTrack.Id]
}

// Crawl each file in the index
func (index *Index) Crawl(callback func(IndexTrack)) {
	for _, v := range index.Tracks {
		callback(*v)
	}
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
