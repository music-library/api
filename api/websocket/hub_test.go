package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	fastWebsocket "github.com/fasthttp/websocket"
	fiberWebsocket "github.com/gofiber/contrib/v3/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHubContinuesAfterPingFailure(t *testing.T) {
	hub := NewHub()
	pings := make(chan time.Time)
	done := make(chan struct{})
	runExited := make(chan struct{})
	connected := make(chan *Client, 2)
	disconnected := make(chan *Client, 2)
	barrier := make(chan bool, 2)

	var failedClient *Client
	hub.On(WsConnect, func(_ *Hub, event *ClientEvent) {
		connected <- event.Client
	})
	hub.On(WsDisconnect, func(_ *Hub, event *ClientEvent) {
		disconnected <- event.Client
	})
	hub.On("test:barrier", func(h *Hub, _ *ClientEvent) {
		_, registered := h.Clients[failedClient]
		barrier <- registered
	})

	go func() {
		hub.run(pings, done)
		close(runExited)
	}()

	t.Cleanup(func() {
		close(done)
		select {
		case <-runExited:
		case <-time.After(time.Second):
			t.Error("hub did not stop")
		}
	})

	failedClient = &Client{
		Hub:  hub,
		Conn: newClosedWebSocketConn(t),
		Id:   "failed-client",
	}
	hub.Register <- failedClient
	assert.Equal(t, failedClient, receiveClient(t, connected, "failed client registration"))

	pings <- time.Now()
	assert.Equal(t, failedClient, receiveClient(t, disconnected, "failed client disconnect"))

	hub.Inbound <- NewClientEvent(failedClient, NewEvent("test:barrier", nil))
	assert.False(t, receiveBool(t, barrier, "inbound event after failed ping"), "failed client should be removed")

	hub.Unregister <- failedClient
	hub.Inbound <- NewClientEvent(failedClient, NewEvent("test:barrier", nil))
	assert.False(t, receiveBool(t, barrier, "inbound event after duplicate unregister"), "failed client should remain removed")
	select {
	case event := <-disconnected:
		t.Fatalf("duplicate unregister emitted a disconnect event for %q", event.Id)
	default:
	}

	followupClient := &Client{Hub: hub, Id: "followup-client"}
	hub.Register <- followupClient
	assert.Equal(t, followupClient, receiveClient(t, connected, "registration after failed ping"))
}

func newClosedWebSocketConn(t *testing.T) *fiberWebsocket.Conn {
	t.Helper()

	clientConn, serverConn := net.Pipe()
	handshakeDone := make(chan error, 1)
	serverDone := make(chan struct{})

	go func() {
		defer close(serverDone)
		defer serverConn.Close()

		request, err := http.ReadRequest(bufio.NewReader(serverConn))
		if err != nil {
			handshakeDone <- err
			return
		}

		challenge := sha1.Sum([]byte(request.Header.Get("Sec-WebSocket-Key") + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
		accept := base64.StdEncoding.EncodeToString(challenge[:])
		_, err = fmt.Fprintf(
			serverConn,
			"HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n",
			accept,
		)
		handshakeDone <- err
		if err == nil {
			_, _ = io.Copy(io.Discard, serverConn)
		}
	}()

	webSocketURL := &url.URL{Scheme: "ws", Host: "example.test", Path: "/"}
	conn, _, err := fastWebsocket.NewClient(clientConn, webSocketURL, nil, 1024, 1024)
	if err != nil {
		t.Fatalf("create websocket connection: %v", err)
	}
	if err := <-handshakeDone; err != nil {
		t.Fatalf("complete websocket handshake: %v", err)
	}
	if err := conn.UnderlyingConn().Close(); err != nil {
		t.Fatalf("close websocket connection: %v", err)
	}
	<-serverDone

	return &fiberWebsocket.Conn{Conn: conn}
}

func receiveClient(t *testing.T, events <-chan *Client, description string) *Client {
	t.Helper()

	select {
	case client := <-events:
		return client
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", description)
		return nil
	}
}

func receiveBool(t *testing.T, events <-chan bool, description string) bool {
	t.Helper()

	select {
	case value := <-events:
		return value
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", description)
		return false
	}
}
