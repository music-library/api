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

	// Register client (connected)
	Register chan *Client

	// Unregister client (disconnect)
	Unregister chan *Client

	// Inbound messages from the clients
	Inbound chan *ClientEvent

	// User defined event handlers for specific incoming events.
	//
	// Use `On` method to register a handler for an event.
	InboundEventHandlers map[string][]func(*Hub, *ClientEvent)
}

func NewHub() *Hub {
	return &Hub{
		Clients:              make(map[*Client]bool),
		Register:             make(chan *Client),
		Unregister:           make(chan *Client),
		Inbound:              make(chan *ClientEvent),
		InboundEventHandlers: make(map[string][]func(*Hub, *ClientEvent)),
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
			h.onCallInboundEventHandlers(NewClientEvent(client, NewEvent(WsConnect, nil)))
			log.WithField("remoteAddr", client.GetIp()).Debug("ws/hub registering client")

			// Unregister client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				log.WithField("remoteAddr", client.GetIp()).WithField("duration", time.Since(time.Unix(client.StartTime, 0))).Debug("ws/hub unregistering client")
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				client.Conn.Close()
				delete(h.Clients, client)
				h.onCallInboundEventHandlers(NewClientEvent(client, NewEvent(WsDisconnect, nil)))
				continue
			}

			// Inbound messages from clients
		case event := <-h.Inbound:
			// Call user defined handlers for this event
			h.onCallInboundEventHandlers(event)

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

// Register handler(s) for an event
func (h *Hub) On(event string, handlers ...func(*Hub, *ClientEvent)) {
	h.InboundEventHandlers[event] = append(h.InboundEventHandlers[event], handlers...)
}

// Internal method to call all handlers for an event
func (h *Hub) onCallInboundEventHandlers(clientEvent *ClientEvent) {
	if handlers, ok := h.InboundEventHandlers[clientEvent.Event.Type]; ok {
		for _, handler := range handlers {
			handler(h, clientEvent)
		}
	}
}

// Emit sends an event to all connected clients.
//
// If clients are specified, the event will be sent only to those clients.
func (h *Hub) Emit(event *Event, clients ...*Client) error {
	msg, err := event.ToJSON()

	if err != nil {
		return err
	}

	log.WithField("wsEvent", event.Type).Debug("ws/hub broadcasting message")

	// Emit to specific clients
	if len(clients) > 0 {
		for _, client := range clients {
			client.Send(msg)
		}

		return nil
	}

	// Emit to all clients
	for client := range h.Clients {
		client.Send(msg)
	}

	return nil
}

func (h *Hub) EmitConnectionCount() error {
	return h.Emit(NewEvent(WsEventConnectionCount, len(h.Clients)))
}
