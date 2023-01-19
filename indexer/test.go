package indexer

import (
	"fmt"

	"github.com/icrowley/fake"
)

func TestGenerateIndex(count uint64) *Index {
	index := Index{
		Name:      "test",
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
