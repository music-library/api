package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	ws "gitlab.com/music-library/music-api/api/websocket"
	"gitlab.com/music-library/music-api/indexer"
)

func TestMusicPlayTrackRejectsInvalidDataTypes(t *testing.T) {
	originalSocket := indexer.MusicLibIndex.Socket
	defer func() {
		indexer.MusicLibIndex.Socket = originalSocket
	}()

	tests := []struct {
		name    string
		payload string
	}{
		{name: "number", payload: `{"type":"music:playTrack","data":123}`},
		{name: "boolean", payload: `{"type":"music:playTrack","data":true}`},
		{name: "array", payload: `{"type":"music:playTrack","data":["track-1"]}`},
		{name: "object", payload: `{"type":"music:playTrack","data":{"id":"track-1"}}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexer.MusicLibIndex.Socket = indexer.NewSocket()
			hub := ws.NewHub()
			WebsocketEventHanders(hub)
			handlers := hub.InboundEventHandlers["music:playTrack"]
			if !assert.Len(t, handlers, 1) {
				t.FailNow()
			}

			existingClient := &ws.Client{Id: "existing-client"}
			existingSession := indexer.MusicLibIndex.Socket.GetOrCreateSession(existingClient.Id)
			existingSession.PlayingTrackId = "existing-track"

			invalidEvent := decodeWebsocketEvent(t, test.payload)
			assert.NotPanics(t, func() {
				handlers[0](hub, ws.NewClientEvent(existingClient, invalidEvent))
			})
			assert.Equal(t, "existing-track", existingSession.PlayingTrackId)

			newClient := &ws.Client{Id: "new-client"}
			assert.NotPanics(t, func() {
				handlers[0](hub, ws.NewClientEvent(newClient, decodeWebsocketEvent(t, test.payload)))
			})
			_, sessionCreated := indexer.MusicLibIndex.Socket.Sessions[newClient.Id]
			assert.False(t, sessionCreated)

			validEvent := decodeWebsocketEvent(t, `{"type":"music:playTrack","data":"next-track"}`)
			handlers[0](hub, ws.NewClientEvent(newClient, validEvent))
			newSession, sessionCreated := indexer.MusicLibIndex.Socket.Sessions[newClient.Id]
			if assert.True(t, sessionCreated) {
				assert.Equal(t, "next-track", newSession.PlayingTrackId)
			}
		})
	}
}

func TestMusicPlayTrackAcceptsStringAndNullData(t *testing.T) {
	originalSocket := indexer.MusicLibIndex.Socket
	defer func() {
		indexer.MusicLibIndex.Socket = originalSocket
	}()

	tests := []struct {
		name            string
		payload         string
		expectedTrackID string
	}{
		{name: "string", payload: `{"type":"music:playTrack","data":"track-1"}`, expectedTrackID: "track-1"},
		{name: "empty string", payload: `{"type":"music:playTrack","data":""}`, expectedTrackID: ""},
		{name: "null", payload: `{"type":"music:playTrack","data":null}`, expectedTrackID: ""},
		{name: "omitted", payload: `{"type":"music:playTrack"}`, expectedTrackID: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexer.MusicLibIndex.Socket = indexer.NewSocket()
			hub := ws.NewHub()
			WebsocketEventHanders(hub)
			handlers := hub.InboundEventHandlers["music:playTrack"]
			if !assert.Len(t, handlers, 1) {
				t.FailNow()
			}

			client := &ws.Client{Id: "test-client"}
			handlers[0](hub, ws.NewClientEvent(client, decodeWebsocketEvent(t, test.payload)))

			session, sessionCreated := indexer.MusicLibIndex.Socket.Sessions[client.Id]
			if assert.True(t, sessionCreated) {
				assert.Equal(t, test.expectedTrackID, session.PlayingTrackId)
			}
		})
	}
}

func decodeWebsocketEvent(t *testing.T, payload string) *ws.Event {
	t.Helper()

	event := &ws.Event{}
	if err := json.Unmarshal([]byte(payload), event); err != nil {
		t.Fatal(err)
	}
	return event
}
