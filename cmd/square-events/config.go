package main

import "os"

type eventServiceConfig struct {
	webhookUrl string
}

func newEventServiceConfig() eventServiceConfig {
	return eventServiceConfig{
		webhookUrl: os.Getenv("WEBHOOK_NOTIFICATION_URL"),
	}
}
