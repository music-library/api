package global

import (
	"github.com/hmerritt/go-ngram"
	"gitlab.com/music-library/music-api/indexer"
)

var IndexMany = indexer.IndexMany{
	DefaultKey: "main",
	Indexes:    make(map[string]*indexer.Index),
}

var Index = indexer.GetNewIndex("main")

var IndexNgram = ngram.NgramIndex{
	NgramMap:   make(map[string]map[int]*ngram.IndexValue),
	IndexesMap: make(map[int]*ngram.IndexValue),
	Ngram:      3,
}
