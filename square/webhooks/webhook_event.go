package webhooks

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
)

const (
	PaymentCreated = "payment.created"
)

func NewWebhookEvent(body string) (WebhookEvent[any], error) {
	var metadata WebhookEventMetadata
	if err := json.Unmarshal([]byte(body), &metadata); err != nil {
		log.Warn().Err(err).Msg("Failed to unmarshal webhook event metadata")
		return nil, fmt.Errorf("malformed request json")
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
