package webhooks

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square/db"
	"github.com/timhugh/digitalvenue/square/queue"
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

func (handler PaymentCreatedHandler) HandleEvent(event WebhookEvent[any]) error {
	paymentCreatedEvent, ok := event.(PaymentCreatedEvent)
	if !ok {
		return fmt.Errorf("event is not PaymentCreatedEvent")
	}
	paymentData, ok := paymentCreatedEvent.Data().(PaymentData)
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

	payment := db.Payment{
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
