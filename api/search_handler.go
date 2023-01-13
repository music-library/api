package api

import (
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/indexer"
)

func SearchHandler(c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "application/json")

	//
	// New behavior (returns object instead of array)
	//
	// tracksJSON, err := sonic.Marshal(global.Index.Tracks)

	//
	// Legacy behavior
	//
	indexValues := global.Ngram.Search(strings.ToLower(c.Params("query")))
	tracksArr := make([]*indexer.IndexTrack, 0, len(indexValues))

	for _, track := range indexValues {
		tracksArr = append(tracksArr, track.Data.(*indexer.IndexTrack))
	}

	tracksJSON, err := sonic.Marshal(tracksArr)

	if err != nil {
		log.Error("http/tracks failed to marshal tracks")
		return Error(c, 500, "failed to marshal tracks")
	}

	return c.Send(tracksJSON)
}
