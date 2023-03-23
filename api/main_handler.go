package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/music-library/music-api/indexer"
)

func MainHandler(c *fiber.Ctx) error {
	libId := c.Locals("libId").(string)
	return c.JSON(indexer.MusicLibIndex.Indexes[libId])
}

func TracksHandler(c *fiber.Ctx) error {
	libId := c.Locals("libId").(string)
	return c.JSON(indexer.MusicLibIndex.Indexes[libId].Tracks)
}

func ReindexHandler(c *fiber.Ctx) error {
	indexer.IndexAllLibraries()
	return c.JSON("ok")
}
