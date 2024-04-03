package sqs

import (
	"github.com/timhugh/digitalvenue/square"
	"testing"
)

func TestSquarePaymentCreatedQueue_SatisfiesInterface(t *testing.T) {
	var _ square.PaymentCreatedQueue = SquarePaymentCreatedQueue{}
}
