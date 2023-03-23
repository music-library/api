package config

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

type Configuration struct {
	// Internal
	LogFile         string    `json:"log_file"`
	LogLevel        string    `json:"log_level"`
	AuthPassword    string    `json:"auth_password"`
	ReIndexHours    string    `json:"re_index_hours"`
	ServerStartTime time.Time `json:"server_start_time"`
	// Music Library
	DataDir      string `json:"data_dir"`
	MusicDir     string `json:"music_dir"`
	MusicDirName string `json:"music_dir_name"`
	//
	MusicDir2      string                      `json:"music_dir2"`
	MusicDir2Name  string                      `json:"music_dir2_name"`
	MusicDir3      string                      `json:"music_dir3"`
	MusicDir3Name  string                      `json:"music_dir3_name"`
	MusicDir4      string                      `json:"music_dir4"`
	MusicDir4Name  string                      `json:"music_dir4_name"`
	MusicDir5      string                      `json:"music_dir5"`
	MusicDir5Name  string                      `json:"music_dir5_name"`
	MusicLibraries []ConfigurationMusicLibrary `json:"music_libraries"`
}

type ConfigurationMusicLibrary struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"-"`
}

var Config = GetConfig() // Global config

func GetConfig() Configuration {
	DATA_DIR := GetEnv("DATA_DIR", "./data")
	MUSIC_DIR := GetEnv("MUSIC_DIR", "./music")

	config := Configuration{
		// Internal
		LogFile:         GetEnv("LOG_FILE", fmt.Sprintf("%s/music-api.log", DATA_DIR)),
		LogLevel:        GetEnv("LOG_LEVEL", "info"),
		AuthPassword:    GetEnv("AUTH_PASSWORD", "lol"),
		ReIndexHours:    GetEnv("REINDEX_HOURS", "12"),
		ServerStartTime: time.Now(),
		// Music Library
		DataDir:      DATA_DIR,
		MusicDir:     MUSIC_DIR,
		MusicDirName: GetEnv("MUSIC_DIR_NAME", "Main"),
		//
		MusicDir2:     GetEnv("MUSIC_DIR2", ""),
		MusicDir2Name: GetEnv("MUSIC_DIR2_NAME", ""),
		MusicDir3:     GetEnv("MUSIC_DIR3", ""),
		MusicDir3Name: GetEnv("MUSIC_DIR3_NAME", ""),
		MusicDir4:     GetEnv("MUSIC_DIR4", ""),
		MusicDir4Name: GetEnv("MUSIC_DIR4_NAME", ""),
		MusicDir5:     GetEnv("MUSIC_DIR5", ""),
		MusicDir5Name: GetEnv("MUSIC_DIR5_NAME", ""),
	}

	config.MusicLibraries = []ConfigurationMusicLibrary{
		{
			Id:   slug.Make(config.MusicDirName),
			Name: config.MusicDirName,
			Path: config.MusicDir,
		},
		{
			Id:   slug.Make(config.MusicDir2Name),
			Name: config.MusicDir2Name,
			Path: config.MusicDir2,
		},
		{
			Id:   slug.Make(config.MusicDir3Name),
			Name: config.MusicDir3Name,
			Path: config.MusicDir3,
		},
		{
			Id:   slug.Make(config.MusicDir4Name),
			Name: config.MusicDir4Name,
			Path: config.MusicDir4,
		},
		{
			Id:   slug.Make(config.MusicDir5Name),
			Name: config.MusicDir5Name,
			Path: config.MusicDir5,
		},
	}

	// Remove empty music library configs
	for i, musicLibConfig := range config.MusicLibraries {
		if len(musicLibConfig.Path) == 0 || len(musicLibConfig.Name) == 0 {
			config.MusicLibraries = config.MusicLibraries[:i]
			break
		}
	}

	return config
}
