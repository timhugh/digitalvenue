package webhooks

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
)

const (
	PaymentCreated = "payment.created"
)

// FIXME: passing the logger feels a little gross -- maybe this should become a class
func NewWebhookEvent(body string, log zerolog.Logger) (WebhookEvent[any], error) {
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

func (base webhookEventBase) EventId() string {
	return base.eventId
}

func (base webhookEventBase) EventType() string {
	return base.eventType
}

func (base webhookEventBase) MerchantId() string {
	return base.merchantId
}
