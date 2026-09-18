// Package realtime is a tiny in-process pub/sub hub used to broadcast seat
// availability changes to connected clients over SSE.
package realtime

import "sync"

// Update is broadcast whenever an event's available seats change.
type Update struct {
	EventID        string `json:"eventId"`
	SeatsAvailable int    `json:"seatsAvailable"`
}

type Hub struct {
	mu          sync.Mutex
	subscribers map[chan Update]struct{}
}

var Default = NewHub()

func NewHub() *Hub {
	return &Hub{subscribers: make(map[chan Update]struct{})}
}

// Subscribe returns a channel that receives updates until Unsubscribe is called.
func (h *Hub) Subscribe() chan Update {
	ch := make(chan Update, 8)
	h.mu.Lock()
	h.subscribers[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch chan Update) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.subscribers[ch]; ok {
		delete(h.subscribers, ch)
		close(ch)
	}
}

// Publish sends an update to every subscriber, dropping it for any subscriber
// whose buffer is full so a slow client cannot block the request path.
func (h *Hub) Publish(update Update) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subscribers {
		select {
		case ch <- update:
		default:
		}
	}
}
