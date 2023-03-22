package indexer

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/icrowley/fake"
	"gitlab.com/music-library/music-api/config"
)

func TestGenerateIndexMany(names []string, count uint64) *IndexMany {
	indexMany := IndexMany{
		DefaultKey: names[0],
		Indexes:    make(map[string]*Index),
	}

	for _, name := range names {
		indexMany.Indexes[name] = TestGenerateIndex(name, count)
	}

	return &indexMany
}

func TestGenerateIndex(name string, count uint64) *Index {
	index := Index{
		Id:        slug.Make(name),
		Name:      slug.Make(name),
		Libraries: config.Config.MusicLibraries,
		Tracks:    make([]*IndexTrack, 0, count),
		TracksKey: make(map[string]int, count),
	}

	for i := 0; i < int(count); i++ {
		metadata := TestGenerateMetadata()
		itemPath := fake.Characters()
		itemId := HashString(itemPath)

		index.TracksKey[itemId] = len(index.Tracks)
		index.Tracks = append(index.Tracks, &IndexTrack{
			Id:       itemId,
			Path:     itemPath,
			Metadata: metadata,
			Stats:    GetEmptyStat(),
		})
	}

	return &index
}

func TestGenerateMetadata() *Metadata {
	return &Metadata{
		Track:       fake.Day(),
		Title:       fake.Characters(),
		Artist:      fake.FullName(),
		AlbumArtist: fake.FullName(),
		Album:       fake.Characters(),
		Year:        fmt.Sprint(fake.Year(1700, 2022)),
		Genre:       fake.Characters(),
		Composer:    fake.FullName(),
		Duration:    fake.Day(),
	}
}
