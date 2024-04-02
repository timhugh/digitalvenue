package sqs

import (
	"github.com/timhugh/digitalvenue/square/queue"
	"testing"
)

func TestPaymentCreatedQueue_SatisfiesInterface(t *testing.T) {
	var _ queue.PaymentCreatedQueue = PaymentCreatedQueue{}
}
