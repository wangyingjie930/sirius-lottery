package eventbus

import (
	"context"
	"sync"
)

// Event represents a generic event.
type Event interface {
	Topic() string
}

// Handler is a function that handles an event.
type Handler func(ctx context.Context, event Event) error

// MemoryEventBus is an in-memory implementation of the EventBus.
type MemoryEventBus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewMemoryEventBus creates a new MemoryEventBus.
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe subscribes to a topic.
func (b *MemoryEventBus) Subscribe(topic string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], handler)
}

// Publish publishes an event to a topic.
func (b *MemoryEventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, ok := b.handlers[event.Topic()]; ok {
		for _, handler := range handlers {
			// Asynchronous execution
			go func(h Handler) {
				_ = h(ctx, event)
			}(handler)
		}
	}
	return nil
}
