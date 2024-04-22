package webhooks

import (
	"context"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
)

type PaymentCreatedHandler struct {
	paymentsRepository square.PaymentRepository
}

func NewPaymentCreatedHandler(paymentsRepository square.PaymentRepository) *PaymentCreatedHandler {
	return &PaymentCreatedHandler{
		paymentsRepository: paymentsRepository,
	}
}

func (handler *PaymentCreatedHandler) HandleEvent(ctx context.Context, event WebhookEvent[any]) error {
	paymentCreatedEvent, ok := event.(*PaymentCreatedEvent)
	if !ok {
		return errors.New("event is not PaymentCreatedEvent")
	}
	paymentData, ok := paymentCreatedEvent.Data().(PaymentData)
	if !ok {
		return errors.New("data type is not PaymentData")
	}

	_, log := logger.FromContext(ctx)
	log.AddParams(map[string]interface{}{
		"squarePaymentID": paymentData.PaymentID,
		"squareOrderID":   paymentData.OrderID,
	})
	log.Info("Processing payment.created event")

	payment := square.Payment{
		SquarePaymentID:  paymentData.PaymentID,
		SquareOrderID:    paymentData.OrderID,
		SquareMerchantID: paymentCreatedEvent.MerchantID(),
	}

	if err := handler.paymentsRepository.PutSquarePayment(&payment); err != nil {
		return errors.Wrap(err, "failed to save payment")
	}

	return nil
}
