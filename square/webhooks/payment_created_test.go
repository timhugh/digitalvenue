package webhooks

import (
	"github.com/matryer/is"
	"testing"
)

func TestPaymentCreatedService_HandleEvent(t *testing.T) {
	t.Skip("TODO: implement")

	is := is.New(t)
	service := PaymentCreatedHandler{}

	event := PaymentCreatedEvent{}
	err := service.HandleEvent(event)

	// TODO

	is.NoErr(err)
}
