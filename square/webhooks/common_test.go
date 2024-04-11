package webhooks

import (
	"github.com/timhugh/digitalvenue/test"
	"os"
)

var paymentCreatedEventRawJson, _ = os.ReadFile("payment-created-event.json")
var paymentCreatedEventJson = string(paymentCreatedEventRawJson)

func newPaymentCreatedEvent() *PaymentCreatedEvent {
	return &PaymentCreatedEvent{
		webhookEventBase: webhookEventBase{
			eventID:    "squareEventID",
			eventType:  "payment.created",
			merchantID: "squareMerchantID",
			tenantID:   test.TenantID,
		},
		data: PaymentData{
			PaymentID: "squarePaymentID",
			OrderID:   "squareOrderID",
		},
	}
}
