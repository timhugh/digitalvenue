package core

import (
	"fmt"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

type EventHandler interface {
	HandleEvent(event webhooks.WebhookEvent[any]) error
}

func GetHandler(eventType string) (EventHandler, error) {
	switch eventType {
	case webhooks.PaymentCreated:
		return PaymentCreatedService{}, nil
	default:
		return nil, fmt.Errorf("unknown event type")
	}
}
