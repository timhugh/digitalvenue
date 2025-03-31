package core

import (
	"log"
	"sync"
)

type Message struct {
	Event   string
	Payload any
}

type MessageQueue struct {
	channels map[string][]chan Message
	mutex    sync.Mutex
}

func New() *MessageQueue {
	return &MessageQueue{
		channels: make(map[string][]chan Message),
	}
}

func (mq *MessageQueue) Subscribe(topic string) chan Message {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	channel := make(chan Message, 100) // TODO: arbitrary buffer size
	mq.channels[topic] = append(mq.channels[topic], channel)
	return channel
}

func (mq *MessageQueue) Publish(topic string, message Message) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	// TODO: write events to database before publishing to channels
	for _, channel := range mq.channels[topic] {
		select {
		case channel <- message:
		default:
			log.Printf("Message dropped for topic %s", topic)
		}
	}
}
