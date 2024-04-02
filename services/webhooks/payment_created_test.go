package webhooks

import (
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/square/webhooks"
	"testing"
)

func TestPaymentCreatedService_HandleEvent(t *testing.T) {
	is := is.New(t)
	service := PaymentCreatedHandler{}

	event := webhooks.PaymentCreatedEvent{}
	err := service.HandleEvent(event)

	// TODO

	is.NoErr(err)
}
