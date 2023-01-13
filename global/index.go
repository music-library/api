package global

import (
	"github.com/hmerritt/go-ngram"
	"gitlab.com/music-library/music-api/indexer"
)

var Index = indexer.Index{
	Tracks:      make(map[string]*indexer.IndexTrack, 5000),
	TracksCount: 0,
}

var Cache = indexer.Cache{
	Path: DATA_DIR,
}

// Initialize the ngrams
var Ngram = ngram.NgramIndex{
	NgramMap:   make(map[string]map[int]*ngram.IndexValue),
	IndexesMap: make(map[int]*ngram.IndexValue),
	Ngram:      3,
}
