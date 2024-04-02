package webhooks

import "os"

var paymentCreatedEventRawJson, _ = os.ReadFile("payment-created-event.json")
var paymentCreatedEventJson = string(paymentCreatedEventRawJson)
