package indexer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type IndexFile struct {
	Id       string
	Path     string
	FilePath string
}

type Index struct {
	Files map[string]*IndexFile
}

func (index *Index) Populate(path string) {
	start := time.Now()
	log.Debug("cache/populate start " + path)

	err := filepath.Walk(path, func(itemPath string, info os.FileInfo, err error) error {

		if err != nil {

			fmt.Println(err)
			return nil
		}

		// @TODO: Check extension
		// filepath.Ext(path) == ".txt"

		if !info.IsDir() {
			fmt.Println(itemPath)

			// @TODO: Add to index
			// files = append(files, itemPath)
		}

		return nil
	})

	if err != nil {
		log.Fatal("cache/populate failed ", err)
	}

	log.Debug(fmt.Sprintf("cache/populate end in %s %s", time.Since(start), path))
}

// Crawl each file in the index
func (index *Index) Crawl(callback func(IndexFile)) {
	for _, v := range index.Files {
		callback(*v)
	}
}
