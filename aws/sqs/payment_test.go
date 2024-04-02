package sqs

import (
	"github.com/timhugh/digitalvenue/square/queue"
	"testing"
)

func TestSquarePaymentCreatedQueue_SatisfiesInterface(t *testing.T) {
	var _ queue.SquarePaymentCreatedQueue = SquarePaymentCreatedQueue{}
}
