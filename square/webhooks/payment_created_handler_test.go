package webhooks

import (
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square/db"
	"github.com/timhugh/digitalvenue/square/queue"
	"testing"
)

func TestPaymentCreatedService_HandleEvent(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentsRepo := mock.Mock[db.PaymentsRepository]()
	paymentCaptor := mock.Captor[db.Payment]()
	mock.WhenSingle(paymentsRepo.CreatePayment(paymentCaptor.Capture())).ThenReturn(nil)

	paymentCreatedQueue := mock.Mock[queue.PaymentCreatedQueue]()
	paymentEventIDCaptor := mock.Captor[string]()
	mock.WhenSingle(paymentCreatedQueue.Publish(paymentEventIDCaptor.Capture())).ThenReturn(nil)

	log := zerolog.Logger{}

	service := PaymentCreatedHandler{
		paymentsRepository:  paymentsRepo,
		paymentCreatedQueue: paymentCreatedQueue,
		log:                 log,
	}

	event := PaymentCreatedEvent{
		webhookEventBase: webhookEventBase{
			eventType:  "payment.created",
			merchantID: "merchant_id",
			eventID:    "event_id",
			body:       paymentCreatedEventJson,
		},
		data: PaymentData{
			PaymentID:  "payment_id",
			LocationID: "location_id",
			OrderID:    "order_id",
		},
	}
	err := service.HandleEvent(event)
	is.NoErr(err)

	is.Equal(paymentCaptor.Last(), db.Payment{
		SquarePaymentID:  "payment_id",
		SquareOrderID:    "order_id",
		SquareMerchantID: "merchant_id",
	})
	is.Equal(paymentEventIDCaptor.Last(), "payment_id")
}
