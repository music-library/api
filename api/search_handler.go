package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"gitlab.com/music-library/music-api/indexer"
)

func SearchHandler(c *fiber.Ctx) error {
	indexValues := indexer.IndexNgram.Search(strings.ToLower(c.Params("query")))
	tracksArr := make([]string, 0, len(indexValues))

	for _, track := range indexValues {
		tracksArr = append(tracksArr, track.Data.(*indexer.IndexTrack).Id)
	}

	return c.JSON(tracksArr)
}
