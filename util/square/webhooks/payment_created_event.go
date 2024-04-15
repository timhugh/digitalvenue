package webhooks

import (
	"encoding/json"
)

type PaymentCreatedEvent struct {
	webhookEventBase
	data PaymentData
}

type PaymentData struct {
	PaymentID string
	OrderID   string
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
					ID      string `json:"id"`
					OrderID string `json:"order_id"`
				} `json:"payment"`
			} `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	event.eventType = raw.EventType
	event.merchantID = raw.MerchantID
	event.eventID = raw.EventID

	event.data.PaymentID = raw.Data.Object.Payment.ID
	event.data.OrderID = raw.Data.Object.Payment.OrderID

	return nil
}
