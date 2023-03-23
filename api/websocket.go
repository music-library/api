package api

import (
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var WsHub = newHub()

func WebsocketUpgradeMiddleware(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func WebsocketHandler(c *websocket.Conn) {
	// c.Locals is added to the *websocket.Conn
	log.Info(c.Locals("allowed"))  // true
	log.Info(c.Query("v"))         // 1.0
	log.Info(c.Cookies("session")) // ""

	client := &Client{hub: WsHub, conn: c, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
