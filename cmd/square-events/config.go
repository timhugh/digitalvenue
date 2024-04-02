package main

import "os"

type eventServiceConfig struct {
	webhookNotificationURL string
}

func newEventServiceConfig() eventServiceConfig {
	return eventServiceConfig{
		webhookNotificationURL: os.Getenv("WEBHOOK_NOTIFICATION_URL"),
	}
}
