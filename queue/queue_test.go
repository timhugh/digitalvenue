package queue_test

import (
	"testing"

	"github.com/timhugh/dv-go/queue"
)

func TestSubscribingAndPublishing(t *testing.T) {
	t.Run("Single subscriber", func(t *testing.T) {
		mq := queue.New()
		topic := "test"
		event := "hello"
		payload := "world"

		channel := mq.Subscribe(topic)

		message := queue.Message{event, payload}
		mq.Publish(topic, message)

		select {
		case msg := <-channel:
			if msg != message {
				t.Errorf("Expected message %q, got %q", message, msg)
			}
		default:
			t.Errorf("Expected message %q, got nothing", message)
		}
	})

	t.Run("Multiple subscribers", func(t *testing.T) {
		mq := queue.New()
		topic := "test"
		event := "hello"
		payload := "world"

		channel1 := mq.Subscribe(topic)
		channel2 := mq.Subscribe(topic)

		message := queue.Message{event, payload}
		mq.Publish(topic, message)

		select {
		case msg := <-channel1:
			if msg != message {
				t.Errorf("Expected message %q, got %q", message, msg)
			}
		default:
			t.Errorf("Expected message %q, got nothing", message)
		}

		select {
		case msg := <-channel2:
			if msg != message {
				t.Errorf("Expected message %q, got %q", message, msg)
			}
		default:
			t.Errorf("Expected message %q, got nothing", message)
		}
	})

	t.Run("Multiple subscribers on different topics", func(t *testing.T) {
		mq := queue.New()
		topic1 := "test1"
		topic2 := "test2"
		event := "hello"
		payload := "world"

		channel1 := mq.Subscribe(topic1)
		channel2 := mq.Subscribe(topic2)

		message1 := queue.Message{event, payload}
		message2 := queue.Message{event, payload}
		mq.Publish(topic1, message1)
		mq.Publish(topic2, message2)

		select {
		case msg := <-channel1:
			if msg != message1 {
				t.Errorf("Expected message %q, got %q", message1, msg)
			}
		default:
			t.Errorf("Expected message %q, got nothing", message1)
		}

		select {
		case msg := <-channel2:
			if msg != message2 {
				t.Errorf("Expected message %q, got %q", message2, msg)
			}
		default:
			t.Errorf("Expected message %q, got nothing", message2)
		}
	})
}
