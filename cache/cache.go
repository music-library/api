package cache

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gitlab.com/music-library/music-api/config"
)

type Cache struct {
	Path string
}

func GetCache(subDir string) Cache {
	path := fmt.Sprintf("%s/%s", config.Config.DataDir, subDir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	return Cache{
		Path: fmt.Sprintf("%s/%s", config.Config.DataDir, subDir),
	}
}

func (cache *Cache) FilePath(path string) string {
	relativePath := cache.Path + "/" + path
	absolutePath, err := filepath.Abs(relativePath)

	if err != nil {
		return relativePath
	}

	return absolutePath
}

func (cache *Cache) Exists(path string) bool {
	fullPath := cache.FilePath(path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Debug("cache/exists file does not exist " + cache.Path + "/" + path)
		return false
	}

	return true
}

func (cache *Cache) Read(path string) ([]byte, error) {
	fullPath := cache.FilePath(path)

	data, err := os.ReadFile(fullPath)

	if err != nil {
		log.Error("cache/read failed to open cache file " + cache.Path + "/" + path)
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
		log.Error("cache/replace failed to add cache file " + cache.Path + "/" + path)
		return err
	}

	log.Debug("cache/replace " + cache.Path + "/" + path)
	return nil
}

func (cache *Cache) Add(path string, fileName string, data []byte) error {
	if !cache.Exists(path + "/" + fileName) {
		return cache.Replace(path, fileName, data)
	}

	return nil
}
