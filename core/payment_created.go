package core

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

type PaymentCreatedService struct {
}

func (s PaymentCreatedService) HandleEvent(event webhooks.WebhookEvent[any]) error {
	paymentCreatedEvent, ok := event.(webhooks.PaymentCreatedEvent)
	if !ok {
		return fmt.Errorf("event is not PaymentCreatedEvent")
	}
	paymentData, ok := paymentCreatedEvent.Data().(webhooks.PaymentData)
	if !ok {
		return fmt.Errorf("data type is not PaymentData")
	}

	log.Debug().
		Str("service", "events-service").
		Str("event", "payment.created").
		Str("merchant_id", paymentCreatedEvent.MerchantId()).
		Str("order_id", paymentData.OrderID).
		Msg("Received event")

	return nil
}
