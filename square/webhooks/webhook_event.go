package webhooks

import (
	"encoding/json"
	"fmt"
)

const (
	PaymentCreated = "payment.created"
)

func NewWebhookEvent(body string) (WebhookEvent[any], error) {
	var metadata WebhookEventMetadata
	if err := json.Unmarshal([]byte(body), &metadata); err != nil {
		return nil, err
	}

	switch metadata.EventType {
	case PaymentCreated:
		var event PaymentCreatedEvent
		if err := json.Unmarshal([]byte(body), &event); err != nil {
			return nil, err
		}
		event.body = body
		return event, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", metadata.EventType)
	}
}

type WebhookEventMetadata struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"type"`
	MerchantID string `json:"merchant_id"`
}

type WebhookEvent[DataType any] interface {
	EventID() string
	EventType() string
	MerchantID() string
	Data() DataType
}

type webhookEventBase struct {
	eventType  string
	merchantID string
	eventID    string
	body       string
}

func (base webhookEventBase) EventID() string {
	return base.eventID
}

func (base webhookEventBase) EventType() string {
	return base.eventType
}

func (base webhookEventBase) MerchantID() string {
	return base.merchantID
}
