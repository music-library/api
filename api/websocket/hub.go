// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package websocket

import (
	"time"

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
			client.StartTime = time.Now().Unix()
			go h.BroadcastConnectionCount()
			log.WithField("remoteAddr", client.GetIp()).Info("ws/hub registering client")

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				log.WithField("remoteAddr", client.GetIp()).WithField("connTime", time.Since(time.Unix(client.StartTime, 0))).Info("ws/hub unregistering client")
				close(client.Send)
				delete(h.Clients, client)
				go h.BroadcastConnectionCount()
			}

		case message := <-h.Broadcast:
			log.Debug("ws/hub broadcasting message", string(message))
			for client := range h.Clients {
				client.Send <- message
			}
		}
	}
}

func (h *Hub) BroadcastConnectionCount() {
	event, _ := NewEvent(WsEventConnectionCount, len(h.Clients)).ToJSON()
	h.Broadcast <- event
}
