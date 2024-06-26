package api

import (
	"github.com/gofiber/fiber/v2"
	fiberWs "github.com/gofiber/websocket/v2"
	"gitlab.com/music-library/music-api/api/websocket"
	"gitlab.com/music-library/music-api/indexer"
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
	h.On(websocket.WsConnect, func(h *websocket.Hub, ce *websocket.ClientEvent) {
		h.EmitConnectionCount()
		EmitPlayingTracks(h, ce.Client) // Only emit playing tracks to the new client
	})

	h.On(websocket.WsDisconnect, func(h *websocket.Hub, ce *websocket.ClientEvent) {
		indexer.MusicLibIndex.Socket.RemoveSession(ce.Client.Id)
		h.EmitConnectionCount()
		EmitPlayingTracks(h)
	})

	h.On("music:playTrack", func(h *websocket.Hub, ce *websocket.ClientEvent) {
		if ce.Event.Data == nil {
			ce.Event.Data = ""
		}

		userId := ce.Client.Id
		playingTrack := ce.Event.Data.(string)

		userSession := indexer.MusicLibIndex.Socket.GetOrCreateSession(userId)
		userSession.PlayingTrackId = playingTrack

		EmitPlayingTracks(h)
	})
}

func EmitPlayingTracks(h *websocket.Hub, clients ...*websocket.Client) {
	h.Emit(websocket.NewEvent("music:playingTracks", indexer.MusicLibIndex.Socket.PlayingTracks()), clients...)
}
