package webhooks

import (
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
)

type PaymentCreatedHandler struct {
	paymentsRepository square.PaymentRepository
	log                *logger.ContextLogger
}

func NewPaymentCreatedHandler(paymentsRepository square.PaymentRepository, log *logger.ContextLogger) *PaymentCreatedHandler {
	return &PaymentCreatedHandler{
		paymentsRepository: paymentsRepository,
		log:                log,
	}
}

func (handler *PaymentCreatedHandler) HandleEvent(event WebhookEvent[any]) error {
	paymentCreatedEvent, ok := event.(*PaymentCreatedEvent)
	if !ok {
		return errors.New("event is not PaymentCreatedEvent")
	}
	paymentData, ok := paymentCreatedEvent.Data().(PaymentData)
	if !ok {
		return errors.New("data type is not PaymentData")
	}

	log := handler.log.Sub().AddParams(map[string]interface{}{
		"event_id":    event.EventID(),
		"payment_id":  paymentData.PaymentID,
		"order_id":    paymentData.OrderID,
		"merchant_id": paymentCreatedEvent.MerchantID(),
		"tenant_id":   event.TenantID(),
	})

	log.Debug("Received event")

	payment := square.Payment{
		SquarePaymentID:  paymentData.PaymentID,
		SquareOrderID:    paymentData.OrderID,
		SquareMerchantID: paymentCreatedEvent.MerchantID(),
	}

	if err := handler.paymentsRepository.PutSquarePayment(&payment); err != nil {
		return errors.Wrap(err, "failed to save payment")
	}

	log.Info("Created payment successfully")

	return nil
}
