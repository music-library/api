package api

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"gitlab.com/music-library/music-api/global"
)

func TracksHandler(c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "application/json")

	trackJSON, err := sonic.Marshal(global.Index.Files)

	if err != nil {
		log.Error("http/tracks failed to marshal tracks")
		return Error(c, 500, "failed to marshal tracks")
	}

	return c.Send(trackJSON)
}
