package webhooks

import (
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/logger"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"testing"
)

func TestPaymentCreatedHandler_HandleEvent(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentsRepo := mock.Mock[square.PaymentRepository]()
	paymentCaptor := mock.Captor[*square.Payment]()
	mock.WhenSingle(paymentsRepo.PutSquarePayment(paymentCaptor.Capture())).ThenReturn(nil)

	service := PaymentCreatedHandler{
		paymentsRepository: paymentsRepo,
		log:                logger.Default(),
	}

	err := service.HandleEvent(newPaymentCreatedEvent())
	is.NoErr(err)

	is.Equal(squaretest.NewSquarePayment(), paymentCaptor.Last())
}
