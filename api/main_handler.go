package api

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"gitlab.com/music-library/music-api/global"
)

func MainHandler(c *fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "application/json")

	indexJSON, err := sonic.Marshal(global.Index)

	if err != nil {
		log.Error("http/tracks failed to marshal index")
		return Error(c, 500, "failed to marshal index")
	}

	return c.Send(indexJSON)
}
