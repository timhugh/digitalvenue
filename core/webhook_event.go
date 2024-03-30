package core

type WebhookEvent struct {
	EventId string `json:"event_id"`
	Type    string `json:"type"`
}
