package square_test

import (
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"github.com/timhugh/digitalvenue/test"
	"testing"
)

func TestMapOrder(t *testing.T) {
	is := is.New(t)

	squareOrder := squaretest.NewSquareOrder()

	order, err := square.MapOrder(squareOrder, squaretest.SquarePaymentID, squaretest.SquareMerchantID)
	is.NoErr(err)

	expectedOrder := squaretest.NewOrder()
	err = test.Diff(order, expectedOrder)
	is.NoErr(err)
}

func TestMapOrder_FailsWithNonIntegerQuantity(t *testing.T) {
	is := is.New(t)

	squareOrder := squaretest.NewSquareOrder()
	squareOrder.OrderItems[0].Quantity = "not an int"

	_, err := square.MapOrder(squareOrder, squaretest.SquarePaymentID, squaretest.SquareMerchantID)
	is.True(err != nil)
}
