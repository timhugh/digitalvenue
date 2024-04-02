package webhooks

import (
	"fmt"
)

type EventHandler interface {
	HandleEvent(event WebhookEvent[any]) error
}

type HandlerProvider struct {
	paymentCreatedHandler PaymentCreatedHandler
}

func NewHandlerProvider(paymentCreatedHandler PaymentCreatedHandler) HandlerProvider {
	return HandlerProvider{
		paymentCreatedHandler: paymentCreatedHandler,
	}
}

func (provider HandlerProvider) GetHandler(eventType string) (EventHandler, error) {
	switch eventType {
	case PaymentCreated:
		return provider.paymentCreatedHandler, nil
	default:
		return nil, fmt.Errorf("unknown event type")
	}
}
