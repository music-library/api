package global

import "gitlab.com/music-library/music-api/indexer"

var Index = indexer.Index{
	Files:      make(map[string]*indexer.IndexFile, 5000),
	FilesCount: 0,
}

var Cache = indexer.Cache{
	Path: DATA_DIR,
}
