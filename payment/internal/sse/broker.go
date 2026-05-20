package sse

import "sync"

// Broker — in-memory pub/sub for SSE.
type Broker struct {
	mu      sync.RWMutex
	clients map[string][]chan string
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[string][]chan string),
	}
}

// Subscribe registers a new SSE client and returns its channel.
func (b *Broker) Subscribe(userID string) chan string {
	ch := make(chan string, 8)
	b.mu.Lock()
	b.clients[userID] = append(b.clients[userID], ch)
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes the client and closes its channel.
func (b *Broker) Unsubscribe(userID string, ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	chans := b.clients[userID]
	for i, c := range chans {
		if c == ch {
			b.clients[userID] = append(chans[:i], chans[i+1:]...)
			close(ch)
			return
		}
	}
}

// Publish broadcasts an event to all connected clients of the user.
func (b *Broker) Publish(userID, payload string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.clients[userID] {
		select {
		case ch <- payload:
		default:
		}
	}
}
