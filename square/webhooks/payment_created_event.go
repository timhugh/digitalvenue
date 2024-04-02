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
	PaymentId  string
	LocationId string
	OrderId    string
}

func (event PaymentCreatedEvent) Data() any {
	return event.data
}

func (event *PaymentCreatedEvent) UnmarshalJSON(data []byte) error {
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

	event.eventType = raw.EventType
	event.merchantId = raw.MerchantId
	event.eventId = raw.EventId

	event.data.PaymentId = raw.Data.Object.Payment.ID
	event.data.LocationId = raw.Data.Object.Payment.LocationId
	event.data.OrderId = raw.Data.Object.Payment.OrderID

	return nil
}
