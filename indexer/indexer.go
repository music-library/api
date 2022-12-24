package indexer

import (
	"fmt"
	"io/ioutil"
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

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("cache/populate failed ", err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	log.Debug(fmt.Sprintf("cache/populate end in %s %s", time.Since(start), path))
}

// Crawl each file in the index
func (index *Index) Crawl(callback func(IndexFile)) {
	for _, v := range index.Files {
		callback(*v)
	}
}
