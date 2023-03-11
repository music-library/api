package global

import (
	"fmt"
	"os"
	"time"
)

var DATA_DIR = GetEnv("DATA_DIR", "./data")
var MUSIC_DIR = GetEnv("MUSIC_DIR", "./music")

// var MUSIC_DIR = GetEnv("MUSIC_DIR", "E:/media/music")

var SERVER_START_TIME = time.Now()

var LOG_LEVEL = GetEnv("LOG_LEVEL", "info")
var LOG_FILE = GetEnv("LOG_FILE", fmt.Sprintf("%s/music-api.log", DATA_DIR))

// Get environment variable or fallback to default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
