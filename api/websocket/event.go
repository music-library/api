package websocket

import (
	"github.com/bytedance/sonic"
)

// Basic event type for every websocket message
type Event struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func NewEvent(event string, data interface{}) *Event {
	return &Event{
		Event: event,
		Data:  data,
	}
}

func (e *Event) Emit(h *Hub) error {
	return h.Emit(e)
}

func (e *Event) ToJSON() ([]byte, error) {
	eventJSON, err := sonic.Marshal(e)

	if err != nil {
		return nil, err
	}

	return eventJSON, nil
}

func (e *Event) ToString() (string, error) {
	eventJSON, err := e.ToJSON()

	if err != nil {
		return "", err
	}

	return string(eventJSON), nil
}

// Internal websocket events
var (
	// Event for when a client connects
	WsEventConnectionCount = "ws:connectionCount"
)
