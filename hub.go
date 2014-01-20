// Package hub provides a simple event dispatcher for publish/subscribe pattern.
package hub

import "sync"

// Event is an interface for published events.
type Event interface {
	Kind() int
}

// Hub is an event dispatcher, publishes events to the subscribers
// which are subscribed for a specific event type.
type Hub struct {
	subscribers map[int][]func(Event)
	sync.RWMutex
}

// New returns pointer to a new Hub.
func New() *Hub {
	return &Hub{subscribers: make(map[int][]func(Event))}
}

// Subscribe for the event of specific kind.
// The caller must receive messages from the retured channel.
// Otherwise, the next Publish() will hang.
func (h *Hub) Subscribe(kind int, handler func(Event)) {
	h.Lock()
	h.subscribers[kind] = append(h.subscribers[kind], handler)
	h.Unlock()
}

// Publish an event to the subscribers.
func (h *Hub) Publish(e Event) {
	h.RLock()
	if handlers, ok := h.subscribers[e.Kind()]; ok {
		for _, handler := range handlers {
			handler(e)
		}
	}
	h.RUnlock()
}

// DefaultHub is the default Hub used by Publish and Subscribe.
var DefaultHub = New()

// Subscribe for the event of specific kind in the DefaultHub.
func Subscribe(kind int, handler func(Event)) {
	DefaultHub.Subscribe(kind, handler)
}

// Publish an event to the subscribers in DefaultHub.
func Publish(e Event) {
	DefaultHub.Publish(e)
}
