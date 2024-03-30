package core

import (
	"fmt"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

const paymentCreated = "payment.created"

type EventHandler interface {
	HandleEvent(event webhooks.WebhookEvent[any]) error
}

func GetHandler(eventType string) (EventHandler, error) {
	switch eventType {
	case paymentCreated:
		return PaymentCreatedService{}, nil
	default:
		return nil, fmt.Errorf("unknown event type")
	}
}
