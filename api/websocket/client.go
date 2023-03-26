// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package websocket

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer. (2kb)
	maxMessageSize = 2000
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Initial client connection time.
	// Useful to calculate connection duration.
	StartTime int64
}

func NewClient(h *Hub, c *websocket.Conn) {
	client := &Client{Hub: h, Conn: c}
	client.Hub.Register <- client

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	client.ReadPump()
}

func (c *Client) Send(message []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	w.Write(message)

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Disconnect() {
	c.Hub.Unregister <- c
}

func (c *Client) GetIp() string {
	return c.Conn.RemoteAddr().String()
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer c.Disconnect()

	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Disconnect()
			}
			break
		}

		c.Hub.Inbound <- message
	}
}
