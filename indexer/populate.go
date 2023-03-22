package indexer

import (
	"sort"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	log "github.com/sirupsen/logrus"
	useCache "gitlab.com/music-library/music-api/cache"
)

// Call this function to index all music libraries
// func IndexAll() {
// 	go (func() {
// 		// Index all music libraries
// 		for _, musicLibConfig := range config.Config.MusicLibraries {
// 			mainIndex := BootstrapIndex(musicLibConfig.Name, musicLibConfig.Path)
// 			global.IndexMany.Indexes[mainIndex.Id] = mainIndex
// 		}
// 	})()
// }

// Async index population (to prevent blocking the server)
func BootstrapIndex(name, dir string) *Index {
	// Populate the index
	newIndex := GetNewIndex(name)
	newIndex.Populate(dir)

	// Read metadata from cache
	cache := useCache.GetCache(newIndex.Id)
	indexCache := ReadAndParseMetadata(cache)

	start := time.Now()
	var await sync.WaitGroup

	// Populate metadata
	for _, indexTrack := range newIndex.Tracks {
		await.Add(1)

		go (func(indexTrack *IndexTrack) {
			defer await.Done()

			// Check if track metadata is cached
			cachedTrackIndex, isCached := indexCache.TracksKey[indexTrack.Id]

			if isCached {
				cachedTrack := indexCache.Tracks[cachedTrackIndex]
				indexTrack.IdAlbum = cachedTrack.IdAlbum
				indexTrack.Metadata = cachedTrack.Metadata
				indexTrack.Stats = cachedTrack.Stats
			} else {
				newIndex.PopulateFileMetadata(indexTrack)
			}

			// Cover
			if !cache.Exists(indexTrack.IdAlbum + "/cover.jpg") {
				trackCover, _ := GetTrackCover(indexTrack.Path)

				if trackCover != nil {
					// Save to global Cache
					cache.Add(indexTrack.IdAlbum, "cover.jpg", trackCover)
					ResizeTrackCover(indexTrack.IdAlbum, "600", cache)
				}
			}
		})(indexTrack)
	}

	await.Wait()

	// Second sync pass
	decadeKeys := make(map[string]bool)
	genresKeys := make(map[string]bool)

	for _, track := range newIndex.Tracks {
		// ngram index
		// global.IndexNgram.Add(GetTrackNgramString(track), ngram.NewIndexValue(index, track))

		// albums
		_, ok := newIndex.Albums[track.IdAlbum]
		if !ok {
			newIndex.Albums[track.IdAlbum] = make([]string, 0, 20)
		}
		newIndex.Albums[track.IdAlbum] = append(newIndex.Albums[track.IdAlbum], track.Id)

		// decades
		if _, ok := decadeKeys[track.Metadata.Decade]; !ok {
			decade := track.Metadata.Decade
			decadeKeys[decade] = true

			if len(decade) == 4 {
				newIndex.Decades = append(newIndex.Decades, decade)
			}
		}

		// genres
		if _, ok := genresKeys[track.Metadata.Genre]; !ok {
			genre := track.Metadata.Genre
			genresKeys[genre] = true

			if len(genre) > 0 {
				newIndex.Genres = append(newIndex.Genres, genre)
			}
		}
	}

	sort.Slice(newIndex.Decades, func(i, j int) bool {
		return newIndex.Decades[i] < newIndex.Decades[j]
	})

	sort.Slice(newIndex.Genres, func(i, j int) bool {
		return newIndex.Genres[i] < newIndex.Genres[j]
	})

	log.Info("main/metadata took ", time.Since(start))

	// Cache metadata
	metadataJSON, err := sonic.Marshal(newIndex)

	if err != nil {
		log.Error("main/metadata/cache failed to marshal metadata ", err)
	}

	cache.Replace(".", "metadata.json", metadataJSON)

	return &newIndex
}
