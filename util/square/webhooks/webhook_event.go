package webhooks

import (
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	PaymentCreated = "payment.created"
)

func NewWebhookEvent(body string) (WebhookEvent[any], error) {
	var metadata WebhookEventMetadata
	if err := json.Unmarshal([]byte(body), &metadata); err != nil {
		return nil, errors.New("failed to unmarshal webhook event metadata")
	}

	switch metadata.EventType {
	case PaymentCreated:
		var event PaymentCreatedEvent
		if err := json.Unmarshal([]byte(body), &event); err != nil {
			return nil, err
		}
		event.body = body
		return &event, nil
	default:
		return nil, errors.Errorf("unknown event type: %s", metadata.EventType)
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
	TenantID() string
	SetTenantID(tenantID string)
	Data() DataType
}

type webhookEventBase struct {
	eventType  string
	merchantID string
	eventID    string
	tenantID   string

	body string
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

func (base webhookEventBase) TenantID() string {
	return base.tenantID
}

func (base *webhookEventBase) SetTenantID(tenantID string) {
	base.tenantID = tenantID
}
