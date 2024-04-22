package webhooks

import (
	"context"
	"errors"
)

type EventHandler interface {
	HandleEvent(ctx context.Context, event WebhookEvent[any]) error
}

type HandlerProvider interface {
	GetHandler(eventType string) (EventHandler, error)
}

type handlerProvider struct {
	paymentCreatedHandler *PaymentCreatedHandler
}

func NewHandlerProvider(paymentCreatedHandler *PaymentCreatedHandler) HandlerProvider {
	return handlerProvider{
		paymentCreatedHandler: paymentCreatedHandler,
	}
}

func (provider handlerProvider) GetHandler(eventType string) (EventHandler, error) {
	switch eventType {
	case PaymentCreated:
		return provider.paymentCreatedHandler, nil
	default:
		return nil, errors.New("unknown event type")
	}
}
