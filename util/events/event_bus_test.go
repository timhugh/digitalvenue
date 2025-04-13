package events_test

import (
	"errors"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/timhugh/digitalvenue/util/events"
)

func TestEventBus(t *testing.T) {
	const bufferSize = 10
	bus := events.NewBus()

	t.Run("emits events to all subscribers", func(t *testing.T) {
		type Event string
		firstSubscription := events.Subscribe[Event](bus, bufferSize)
		defer firstSubscription.Unsubscribe()
		secondSubscription := events.Subscribe[Event](bus, bufferSize)
		defer secondSubscription.Unsubscribe()

		firstDone := make(chan struct{})
		secondDone := make(chan struct{})

		firstSubEvents := make([]Event, 0)
		go func() {
			for event := range firstSubscription.Events {
				firstSubEvents = append(firstSubEvents, event)
				if len(firstSubEvents) == 2 {
					close(firstDone)
				}
			}
		}()

		secondSubEvents := make([]Event, 0)
		go func() {
			for event := range secondSubscription.Events {
				secondSubEvents = append(secondSubEvents, event)
				if len(secondSubEvents) == 2 {
					close(secondDone)
				}
			}
		}()

		events.Emit[Event](bus, "Hello!")
		events.Emit[Event](bus, "Goodbye!")

		select {
		case <-firstDone:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for first subscriber to receive events")
		}

		select {
		case <-secondDone:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for second subscriber to receive events")
		}

		firstSubscription.Unsubscribe()
		secondSubscription.Unsubscribe()

		expectedEvents := []Event{"Hello!", "Goodbye!"}
		firstDiff := deep.Equal(firstSubEvents, expectedEvents)
		if firstDiff != nil {
			t.Errorf("first subscription events don't match expected:\n%v", firstDiff)
		}

		secondDiff := deep.Equal(secondSubEvents, expectedEvents)
		if secondDiff != nil {
			t.Errorf("second subscription events don't match expected:\n%v", secondDiff)
		}
	})

	t.Run("handles multiple event types simultaneously", func(t *testing.T) {
		type StringEvent string
		type StructEvent struct {
			Id      int
			Message string
		}

		stringSubscription := events.Subscribe[StringEvent](bus, bufferSize)
		structSubscription := events.Subscribe[StructEvent](bus, bufferSize)

		stringSubDone := make(chan struct{})
		structSubDone := make(chan struct{})

		stringEvents := make([]StringEvent, 0)
		go func() {
			for event := range stringSubscription.Events {
				stringEvents = append(stringEvents, event)
				if len(stringEvents) == 2 {
					close(stringSubDone)
				}
			}
		}()

		structEvents := make([]StructEvent, 0)
		go func() {
			for event := range structSubscription.Events {
				structEvents = append(structEvents, event)
				if len(structEvents) == 2 {
					close(structSubDone)
				}
			}
		}()

		events.Emit[StringEvent](bus, "Hello!")
		events.Emit[StructEvent](bus, StructEvent{27, "Hello!"})
		events.Emit[StringEvent](bus, "Goodbye!")
		events.Emit[StructEvent](bus, StructEvent{42, "Goodbye!"})

		select {
		case <-stringSubDone:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for string subscriber to receive events")
		}

		select {
		case <-structSubDone:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for struct subscriber to receive events")
		}

		stringSubscription.Unsubscribe()
		structSubscription.Unsubscribe()

		expectedStringEvents := []StringEvent{"Hello!", "Goodbye!"}
		stringDiff := deep.Equal(stringEvents, expectedStringEvents)
		if stringDiff != nil {
			t.Errorf("string subscription events don't match expected:\n%v", stringDiff)
		}

		expectedStructEvents := []StructEvent{
			{27, "Hello!"},
			{42, "Goodbye!"},
		}
		structDiff := deep.Equal(structEvents, expectedStructEvents)
		if structDiff != nil {
			t.Errorf("struct subscription events don't match expected:\n%v", structDiff)
		}
	})

	t.Run("returns errors when queues are full", func(t *testing.T) {
		type Event struct{}
		subscription := events.Subscribe[Event](bus, 1)

		events.Emit[Event](bus, Event{})
		errs := events.Emit[Event](bus, Event{})
		if len(errs) != 1 {
			t.Errorf("expected 1 error, got %d", len(errs))
		}
		if err := errs[0]; !errors.Is(err, events.ErrQueueFull) {
			t.Errorf("expected ErrQueueFull, got %v", err)
		}

		subscription.Unsubscribe()
	})
}
