package global

import (
	"github.com/hmerritt/go-ngram"
	"gitlab.com/music-library/music-api/indexer"
)

var IndexMany = indexer.IndexMany{
	DefaultKey: "Main",
	Indexes:    make(map[string]*indexer.Index),
}

var Index = indexer.GetNewIndex("main")

var Cache = indexer.Cache{
	Path: DATA_DIR,
}

var IndexNgram = ngram.NgramIndex{
	NgramMap:   make(map[string]map[int]*ngram.IndexValue),
	IndexesMap: make(map[int]*ngram.IndexValue),
	Ngram:      3,
}
