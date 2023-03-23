// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package api

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			go h.BroadcastConnectionCount()
			client.remoteAddr = client.conn.RemoteAddr().String()
			log.Debug("ws/hub registering client", client.remoteAddr)
		case client := <-h.unregister:
			log.Debug("ws/hub unregistering client", client.remoteAddr)
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			go h.BroadcastConnectionCount()
		case message := <-h.broadcast:
			log.Debug("ws/hub broadcasting message", string(message))
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastConnectionCount() {
	h.broadcast <- []byte(fmt.Sprintf("There are %d clients connected", len(h.clients)))
}
