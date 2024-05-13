package webhooks

import (
	"context"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/util/square"
	"github.com/timhugh/digitalvenue/util/square/squaretest"
	"testing"
)

func TestPaymentCreatedHandler_HandleEvent(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentsRepo := mock.Mock[square.PaymentRepository]()
	paymentsRepoPaymentCaptor := mock.Captor[*square.Payment]()
	mock.WhenSingle(paymentsRepo.PutSquarePayment(paymentsRepoPaymentCaptor.Capture())).ThenReturn(nil)

	paymentCreatedQueue := mock.Mock[square.PaymentCreatedQueue]()
	paymentCreatedQueuePaymentCaptor := mock.Captor[*square.Payment]()
	mock.When(paymentCreatedQueue.PublishPaymentCreated(paymentCreatedQueuePaymentCaptor.Capture())).ThenReturn(nil)

	service := NewPaymentCreatedHandler(paymentsRepo, paymentCreatedQueue)

	err := service.HandleEvent(context.Background(), newPaymentCreatedEvent())
	is.NoErr(err)

	expectedPayment := squaretest.NewSquarePayment()
	is.Equal(expectedPayment, paymentsRepoPaymentCaptor.Last())
	is.Equal(expectedPayment, paymentCreatedQueuePaymentCaptor.Last())
}
