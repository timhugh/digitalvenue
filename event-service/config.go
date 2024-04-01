package main

import "os"

type EventServiceConfig struct {
	WebhookUrl string
}

func NewEventServiceConfig() EventServiceConfig {
	return EventServiceConfig{
		WebhookUrl: os.Getenv("WEBHOOK_NOTIFICATION_URL"),
	}
}
