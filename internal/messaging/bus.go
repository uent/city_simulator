package messaging

import (
	"context"
	"fmt"
	"sync"
)

// Actor is the interface every simulation participant must implement.
type Actor interface {
	ID() string
	Inbox() chan Message
	Start(ctx context.Context)
	Stop()
}

// MessageBus routes messages to registered actors by ID.
type MessageBus struct {
	mu     sync.RWMutex
	actors map[string]Actor
}

// NewMessageBus returns an empty MessageBus.
func NewMessageBus() *MessageBus {
	return &MessageBus{actors: make(map[string]Actor)}
}

// Register adds an actor to the bus. A second call with the same ID overwrites.
func (b *MessageBus) Register(a Actor) {
	b.mu.Lock()
	b.actors[a.ID()] = a
	b.mu.Unlock()
}

// Send delivers msg to the actor identified by msg.To.
// Returns an error if the actor is not registered.
func (b *MessageBus) Send(msg Message) error {
	b.mu.RLock()
	a, ok := b.actors[msg.To]
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("messaging: actor %q not registered", msg.To)
	}
	a.Inbox() <- msg
	return nil
}

// Broadcast delivers a copy of msg (each with its own ReplyChan) to every registered
// actor except those in excludeIDs. Returns a receive-only fan-in channel that emits
// one reply per recipient and is closed after all replies are forwarded.
// Returns an error if there are no eligible recipients.
func (b *MessageBus) Broadcast(msg Message, excludeIDs ...string) (<-chan Message, error) {
	b.mu.RLock()
	excluded := make(map[string]bool, len(excludeIDs))
	for _, id := range excludeIDs {
		excluded[id] = true
	}
	var targets []Actor
	for id, a := range b.actors {
		if !excluded[id] {
			targets = append(targets, a)
		}
	}
	b.mu.RUnlock()

	if len(targets) == 0 {
		return nil, fmt.Errorf("messaging: no eligible recipients for broadcast")
	}

	fanIn := make(chan Message, len(targets))
	var wg sync.WaitGroup
	for _, a := range targets {
		wg.Add(1)
		go func(actor Actor) {
			defer wg.Done()
			m := msg
			m.ReplyChan = make(chan Message, 1)
			actor.Inbox() <- m
			fanIn <- <-m.ReplyChan
		}(a)
	}
	go func() {
		wg.Wait()
		close(fanIn)
	}()

	return fanIn, nil
}

// StartAll calls Start on every registered actor.
func (b *MessageBus) StartAll(ctx context.Context) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, a := range b.actors {
		a.Start(ctx)
	}
}

// Count returns the number of registered actors.
func (b *MessageBus) Count() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.actors)
}
