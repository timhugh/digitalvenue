package webhooks

import (
	"encoding/json"
	"fmt"
)

func NewWebhookEvent(body string) (WebhookEvent[any], error) {
	var metadata WebhookEventMetadata
	if err := json.Unmarshal([]byte(body), &metadata); err != nil {
		return nil, fmt.Errorf("error unmarshalling webhook event metadata: %w", err)
	}

	switch metadata.EventType {
	case "payment.created":
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
	EventId    string `json:"event_id"`
	EventType  string `json:"type"`
	MerchantId string `json:"merchant_id"`
}

type WebhookEvent[DataType any] interface {
	EventId() string
	EventType() string
	MerchantId() string
	Data() DataType
}

type webhookEventBase struct {
	eventType  string
	merchantId string
	eventId    string
	body       string
}

func (w webhookEventBase) EventId() string {
	return w.eventId
}

func (w webhookEventBase) EventType() string {
	return w.eventType
}

func (w webhookEventBase) MerchantId() string {
	return w.merchantId
}
