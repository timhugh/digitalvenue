package sqs

import (
	"github.com/timhugh/digitalvenue/queue"
	"testing"
)

func PaymentCreatedQueue_SatisfiesInterface(t *testing.T) {
	var _ queue.PaymentCreatedQueue = PaymentCreatedQueue{}
}
