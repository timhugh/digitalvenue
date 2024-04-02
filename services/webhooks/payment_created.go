package webhooks

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/db"
	"github.com/timhugh/digitalvenue/queue"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

type PaymentCreatedHandler struct {
	paymentsRepository  db.PaymentsRepository
	paymentCreatedQueue queue.PaymentCreatedQueue
	log                 zerolog.Logger
}

func NewPaymentCreatedService(paymentsRepository db.PaymentsRepository, paymentCreatedQueue queue.PaymentCreatedQueue, log zerolog.Logger) PaymentCreatedHandler {
	return PaymentCreatedHandler{
		paymentsRepository:  paymentsRepository,
		paymentCreatedQueue: paymentCreatedQueue,
		log:                 log,
	}
}

func (handler PaymentCreatedHandler) HandleEvent(event webhooks.WebhookEvent[any]) error {
	paymentCreatedEvent, ok := event.(webhooks.PaymentCreatedEvent)
	if !ok {
		return fmt.Errorf("event is not PaymentCreatedEvent")
	}
	paymentData, ok := paymentCreatedEvent.Data().(webhooks.PaymentData)
	if !ok {
		return fmt.Errorf("data type is not PaymentData")
	}

	handler.log.Debug().
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

	err := handler.paymentsRepository.CreatePayment(payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	if err := handler.paymentCreatedQueue.Publish(payment.SquarePaymentID); err != nil {
		return fmt.Errorf("failed to publish payment created event: %w", err)
	}

	return nil
}
