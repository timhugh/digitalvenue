package webhooks

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square"
)

type PaymentCreatedHandler struct {
	paymentsRepository  square.PaymentRepository
	paymentCreatedQueue square.PaymentCreatedQueue
	log                 zerolog.Logger
}

func NewPaymentCreatedHandler(paymentsRepository square.PaymentRepository, paymentCreatedQueue square.PaymentCreatedQueue, log zerolog.Logger) *PaymentCreatedHandler {
	return &PaymentCreatedHandler{
		paymentsRepository:  paymentsRepository,
		paymentCreatedQueue: paymentCreatedQueue,
		log:                 log,
	}
}

func (handler *PaymentCreatedHandler) HandleEvent(event WebhookEvent[any]) error {
	paymentCreatedEvent, ok := event.(PaymentCreatedEvent)
	if !ok {
		return errors.New("event is not PaymentCreatedEvent")
	}
	paymentData, ok := paymentCreatedEvent.Data().(PaymentData)
	if !ok {
		return errors.New("data type is not PaymentData")
	}

	handler.log.Debug().
		Str("service", "events-service").
		Str("event", "payment.created").
		Str("payment_id", paymentData.PaymentID).
		Str("order_id", paymentData.OrderID).
		Str("merchant_id", paymentCreatedEvent.MerchantID()).
		Msg("Received event")

	payment := square.Payment{
		SquarePaymentID:  paymentData.PaymentID,
		SquareOrderID:    paymentData.OrderID,
		SquareMerchantID: paymentCreatedEvent.MerchantID(),
	}

	if err := handler.paymentsRepository.PutSquarePayment(payment); err != nil {
		return errors.Wrap(err, "failed to save payment")
	}

	if err := handler.paymentCreatedQueue.Publish(payment.SquarePaymentID); err != nil {
		return errors.Wrap(err, "failed to publish payment created event")
	}

	return nil
}
