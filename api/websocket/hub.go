// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package websocket

import (
	"time"

	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Connected clients
	Clients map[*Client]bool

	// Inbound messages from the clients
	Inbound chan []byte

	// Register client (connected)
	Register chan *Client

	// Unregister client (disconnect)
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Inbound:    make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		// Register client
		case client := <-h.Register:
			h.Clients[client] = true
			client.StartTime = time.Now().Unix()
			go h.EmitConnectionCount()
			log.WithField("remoteAddr", client.GetIp()).Info("ws/hub registering client")

			// Unregister client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				log.WithField("remoteAddr", client.GetIp()).WithField("duration", time.Since(time.Unix(client.StartTime, 0))).Info("ws/hub unregistering client")
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				client.Conn.Close()
				delete(h.Clients, client)
				go h.EmitConnectionCount()
			}

			// Inbound messages from clients
		case message := <-h.Inbound:
			log.Debug("ws/hub incomming message", string(message))
			h.Emit(NewEvent("ws:inbound", string(message)))

			// Ping all clients periodically to check if they are still connected.
			// Disconnect them if they do not respond before the `writeWait` timeout.
		case <-ticker.C:
			for client := range h.Clients {
				client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					client.Disconnect()
				}
			}
		}

	}
}

// Emit sends an event to all connected clients
// @TODO: Emit to specific clients
func (h *Hub) Emit(event *Event) error {
	msg, err := event.ToJSON()

	if err != nil {
		return err
	}

	log.WithField("wsMsg", string(msg)).Debug("ws/hub broadcasting message")

	for client := range h.Clients {
		client.Send(msg)
	}

	return nil
}

func (h *Hub) EmitConnectionCount() error {
	return h.Emit(NewEvent(WsEventConnectionCount, len(h.Clients)))
}
