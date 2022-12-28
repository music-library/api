package api

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"gitlab.com/music-library/music-api/global"
	"gitlab.com/music-library/music-api/indexer"
)

func TracksHandler(c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "application/json")

	//
	// New behavior (returns object instead of array)
	//
	// tracksJSON, err := sonic.Marshal(global.Index.Tracks)

	//
	// Legacy behavior
	//
	tracksArr := make([]*indexer.IndexTrack, 0, len(global.Index.Tracks))
	for _, track := range global.Index.Tracks {
		tracksArr = append(tracksArr, track)
	}

	tracksJSON, err := sonic.Marshal(tracksArr)

	if err != nil {
		log.Error("http/tracks failed to marshal tracks")
		return Error(c, 500, "failed to marshal tracks")
	}

	return c.Send(tracksJSON)
}
