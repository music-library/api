package api

import (
	"github.com/gofiber/fiber/v2"

	"gitlab.com/music-library/music-api/global"
)

func MainHandler(c *fiber.Ctx) error {
	libId := c.Locals("libId").(string)
	return c.JSON(global.IndexMany.Indexes[libId])
}

func TracksHandler(c *fiber.Ctx) error {
	libId := c.Locals("libId").(string)
	return c.JSON(global.IndexMany.Indexes[libId].Tracks)
}
