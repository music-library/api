package api

import (
	"github.com/gofiber/fiber/v2"

	"gitlab.com/music-library/music-api/global"
)

func MainHandler(c *fiber.Ctx) error {
	return c.JSON(global.IndexMany.Indexes[global.IndexMany.DefaultKey])
	// return c.JSON(global.Index)
}

func TracksHandler(c *fiber.Ctx) error {
	return c.JSON(global.Index.Tracks)
}
