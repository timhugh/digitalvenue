package webhooks

import (
	"github.com/matryer/is"
	"github.com/rs/zerolog"
	"testing"
)

func TestPaymentCreatedEvent_Unmarshal(t *testing.T) {
	is := is.New(t)

	log := zerolog.Logger{}
	event, err := NewWebhookEvent(paymentCreatedEventJson, log)
	is.NoErr(err)

	paymentCreatedEvent, ok := event.(PaymentCreatedEvent)
	is.True(ok)

	is.Equal(paymentCreatedEvent.EventID(), "event_id")
	is.Equal(paymentCreatedEvent.EventType(), "payment.created")
	is.Equal(paymentCreatedEvent.MerchantID(), "merchant_id")

	data := paymentCreatedEvent.Data()
	paymentData, ok := data.(PaymentData)
	is.True(ok)
	is.Equal(paymentData.PaymentID, "payment_id")
	is.Equal(paymentData.LocationID, "location_id")
	is.Equal(paymentData.OrderID, "order_id")
}
