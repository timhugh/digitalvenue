package webhooks

import (
	"encoding/json"
	"fmt"
)

type PaymentCreatedEvent struct {
	webhookEventBase
	data PaymentData
}

type PaymentData struct {
	PaymentID  string
	LocationID string
	OrderID    string
}

func (p PaymentCreatedEvent) Data() any {
	return p.data
}

func (p *PaymentCreatedEvent) UnmarshalJSON(data []byte) error {
	var raw struct {
		WebhookEventMetadata
		Data struct {
			Object struct {
				Payment struct {
					ID         string `json:"id"`
					LocationId string `json:"location_id"`
					OrderID    string `json:"order_id"`
				} `json:"payment"`
			} `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("error unmarshalling payment created event: %w", err)
	}

	p.eventType = raw.EventType
	p.merchantId = raw.MerchantId
	p.eventId = raw.EventId

	p.data.PaymentID = raw.Data.Object.Payment.ID
	p.data.LocationID = raw.Data.Object.Payment.LocationId
	p.data.OrderID = raw.Data.Object.Payment.OrderID

	return nil
}
