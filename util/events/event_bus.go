package events

import (
	"fmt"
	"sync"
)

var (
	ErrNoSubscriptions = fmt.Errorf("no subscribers for event")
	ErrQueueFull       = fmt.Errorf("subscriber queue is full")
)

type Bus struct {
	idCounter     SubscriptionId
	subscriptions map[string]map[SubscriptionId]any
	mu            sync.RWMutex
}

type SubscriptionId uint64

type Subscription[E any] struct {
	ID     SubscriptionId
	Events <-chan E
	bus    *Bus
}

func NewBus() *Bus {
	return &Bus{
		subscriptions: make(map[string]map[SubscriptionId]any),
		idCounter:     0,
	}
}

func Subscribe[E any](bus *Bus, bufferSize int) Subscription[E] {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	eventType := getTypeName[E]()

	_, exists := bus.subscriptions[eventType]
	if !exists {
		bus.subscriptions[eventType] = make(map[SubscriptionId]any)
	}

	id := bus.idCounter
	bus.idCounter++

	ch := make(chan E, bufferSize)

	bus.subscriptions[eventType][id] = ch

	return Subscription[E]{
		ID:     id,
		Events: ch,
		bus:    bus,
	}
}

func (s Subscription[E]) Unsubscribe() {
	s.bus.mu.Lock()
	defer s.bus.mu.Unlock()

	typeName := getTypeName[E]()
	subscriptions, exists := s.bus.subscriptions[typeName]
	if !exists {
		return
	}

	ch, found := subscriptions[s.ID]
	if !found {
		return
	}

	delete(subscriptions, s.ID)
	if channel, ok := ch.(chan E); ok {
		close(channel)
	}
}

func Emit[E any](bus *Bus, event E) []error {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	typeName := getTypeName[E]()
	subscriptions, exists := bus.subscriptions[typeName]
	if !exists {
		return []error{ErrNoSubscriptions}
	}

	var errors []error
	for _, subscriber := range subscriptions {
		ch, ok := subscriber.(chan E)
		if !ok {
			continue
		}

		select {
		case ch <- event:
		default:
			errors = append(errors, ErrQueueFull)
		}
	}

	return errors
}

func getTypeName[T any]() string {
	var zero T
	return fmt.Sprintf("%T", zero)
}
