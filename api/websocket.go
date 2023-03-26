package api

import (
	"github.com/gofiber/fiber/v2"
	fiberWs "github.com/gofiber/websocket/v2"
	"gitlab.com/music-library/music-api/api/websocket"
)

var WsHub = websocket.NewHub()

func WebsocketUpgradeMiddleware(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if fiberWs.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func WebsocketHandler(c *fiberWs.Conn) {
	websocket.NewClient(WsHub, c)
}

func WebsocketEventHanders(h *websocket.Hub) {
	h.On("ws:connect", func(h *websocket.Hub, ce *websocket.ClientEvent) {
		h.EmitConnectionCount()
	})

	h.On("ws:disconnect", func(h *websocket.Hub, ce *websocket.ClientEvent) {
		h.EmitConnectionCount()
	})
}
