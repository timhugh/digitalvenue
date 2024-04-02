package webhooks

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/db"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

type PaymentCreatedHandler struct {
	paymentsRepository db.PaymentsRepository
}

func NewPaymentCreatedService(paymentsRepository db.PaymentsRepository) PaymentCreatedHandler {
	return PaymentCreatedHandler{
		paymentsRepository: paymentsRepository,
	}
}

func (s PaymentCreatedHandler) HandleEvent(event webhooks.WebhookEvent[any]) error {
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
		Str("payment_id", paymentData.PaymentID).
		Str("order_id", paymentData.OrderID).
		Str("merchant_id", paymentCreatedEvent.MerchantId()).
		Msg("Received event")

	payment := core.Payment{
		SquarePaymentID:  paymentData.PaymentID,
		SquareOrderID:    paymentData.OrderID,
		SquareMerchantID: paymentCreatedEvent.MerchantId(),
	}

	err := s.paymentsRepository.CreatePayment(payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}
