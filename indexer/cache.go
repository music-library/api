package indexer

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Cache struct {
	Path string
}

func (cache *Cache) FilePath(path string) string {
	return cache.Path + "/" + path
}

func (cache *Cache) Exists(path string) bool {
	fullPath := cache.FilePath(path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Debug("cache/exists file does not exist " + fullPath)
		return false
	}

	return true
}

func (cache *Cache) Read(path string) ([]byte, error) {
	fullPath := cache.FilePath(path)

	data, err := os.ReadFile(fullPath)

	if err != nil {
		log.Error("cache/read failed to open cache file " + fullPath)
		return nil, err
	}

	return data, nil
}

func (cache *Cache) Replace(path string, fileName string, data []byte) error {
	fullPath := cache.FilePath(path + "/" + fileName)

	if _, err := os.Stat(cache.Path + "/" + path); os.IsNotExist(err) {
		os.Mkdir(cache.Path+"/"+path, 0755)
	}

	err := os.WriteFile(fullPath, data, 0766)

	if err != nil {
		log.Error("cache/replace failed to add cache file " + fullPath)
		return err
	}

	log.Debug("cache/replace " + fullPath)
	return nil
}

func (cache *Cache) Add(path string, fileName string, data []byte) error {
	if !cache.Exists(path + "/" + fileName) {
		return cache.Replace(path, fileName, data)
	}

	return nil
}
