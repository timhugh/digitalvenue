package webhooks

import (
	"fmt"
)

type EventHandler interface {
	HandleEvent(event WebhookEvent[any]) error
}

type HandlerProvider interface {
	GetHandler(eventType string) (EventHandler, error)
}

type DefaultHandlerProvider struct {
	paymentCreatedHandler PaymentCreatedHandler
}

func NewHandlerProvider(paymentCreatedHandler PaymentCreatedHandler) HandlerProvider {
	return DefaultHandlerProvider{
		paymentCreatedHandler: paymentCreatedHandler,
	}
}

func (provider DefaultHandlerProvider) GetHandler(eventType string) (EventHandler, error) {
	switch eventType {
	case PaymentCreated:
		return provider.paymentCreatedHandler, nil
	default:
		return nil, fmt.Errorf("unknown event type")
	}
}
