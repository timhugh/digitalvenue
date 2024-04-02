package webhooks

import (
	"fmt"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

type EventHandler interface {
	HandleEvent(event webhooks.WebhookEvent[any]) error
}

type HandlerProvider struct {
	paymentCreatedHandler PaymentCreatedHandler
}

func NewHandlerProvider(paymentCreatedHandler PaymentCreatedHandler) HandlerProvider {
	return HandlerProvider{
		paymentCreatedHandler: paymentCreatedHandler,
	}
}

func (p HandlerProvider) GetHandler(eventType string) (EventHandler, error) {
	switch eventType {
	case webhooks.PaymentCreated:
		return p.paymentCreatedHandler, nil
	default:
		return nil, fmt.Errorf("unknown event type")
	}
}
