package indexer

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type IndexFile struct {
	Id       string
	Path     string
	FileName string
	// Metadata *Metadata
}

type Index struct {
	Files  map[string]*IndexFile
	Albums [][]string // Slice of IndexFile.Id
}

// Populate File index with audio `IndexFile` objects
func (index *Index) Populate(path string, populateMetadata bool) {
	start := time.Now()
	log.Debug("cache/populate start " + path)

	err := filepath.Walk(path, func(itemPath string, info os.FileInfo, err error) error {

		if err != nil {

			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && isFileAudio(itemPath) {
			itemId := hashString(itemPath)

			index.Files[itemId] = &IndexFile{
				Id:       itemId,
				Path:     itemPath,
				FileName: filepath.Base(itemPath), // @Investigate: Is this needed?
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal("cache/populate failed ", err)
	}

	log.Debug(fmt.Sprintf("cache/populate end. indexed %d items in %s %s", len(index.Files), time.Since(start), path))
}

// Crawl each file in the index
func (index *Index) Crawl(callback func(IndexFile)) {
	for _, v := range index.Files {
		callback(*v)
	}
}

// Check if file is audio file
func isFileAudio(path string) bool {
	return filepath.Ext(path) == ".flac" || filepath.Ext(path) == ".mp3" || filepath.Ext(path) == ".ogg"
}

// Hash string using SHA1, returns a string.
func hashString(str string) string {
	// Store string in buffer array
	buf := []byte(str)

	// Hash string, returns buffer array
	hash_byte := sha1.Sum(buf)

	// Convert buffer into hex string
	hash := hex.EncodeToString(hash_byte[:])

	// Return hash as string
	return hash
}
