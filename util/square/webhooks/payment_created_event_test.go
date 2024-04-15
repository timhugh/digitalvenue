package webhooks

import (
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func TestPaymentCreatedEvent_Unmarshal(t *testing.T) {
	is := is.New(t)

	event, err := NewWebhookEvent(paymentCreatedEventJson)
	is.NoErr(err)

	paymentCreatedEvent, ok := event.(*PaymentCreatedEvent)
	is.True(ok)

	expectedEvent := newPaymentCreatedEvent()

	if err := test.Diff(expectedEvent, paymentCreatedEvent); err != nil {
		t.Fatal(err)
	}
}
