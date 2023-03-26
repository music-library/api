// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package websocket

import (
	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			go h.BroadcastConnectionCount()
			client.RemoteAddr = client.Conn.RemoteAddr().String()
			log.Debug("ws/hub registering client", client.RemoteAddr)
		case client := <-h.Unregister:
			log.Debug("ws/hub unregistering client", client.RemoteAddr)
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			go h.BroadcastConnectionCount()
		case message := <-h.Broadcast:
			log.Debug("ws/hub broadcasting message", string(message))
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastConnectionCount() {
	event, _ := NewEvent("_ws:connCount", len(h.Clients)).ToJSON()
	h.Broadcast <- event
}
