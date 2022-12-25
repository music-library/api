package global

import (
	"gitlab.com/music-library/music-api/indexer"
)

var Index = indexer.Index{
	Files: make(map[string]*indexer.IndexFile, 1000),
}
