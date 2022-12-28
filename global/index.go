package global

import "gitlab.com/music-library/music-api/indexer"

var Index = indexer.Index{
	Tracks:      make(map[string]*indexer.IndexTrack, 5000),
	TracksCount: 0,
}

var Cache = indexer.Cache{
	Path: DATA_DIR,
}
