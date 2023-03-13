package global

import (
	"github.com/hmerritt/go-ngram"
	"gitlab.com/music-library/music-api/indexer"
)

var Index = indexer.GetNewIndex("main")

var Cache = indexer.Cache{
	Path: DATA_DIR,
}

// Initialize the ngrams
var IndexNgram = ngram.NgramIndex{
	NgramMap:   make(map[string]map[int]*ngram.IndexValue),
	IndexesMap: make(map[int]*ngram.IndexValue),
	Ngram:      3,
}
