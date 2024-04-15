package square_test

import (
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/util/square"
	"github.com/timhugh/digitalvenue/util/square/squaretest"
	"github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func TestMapOrder(t *testing.T) {
	is := is.New(t)

	squareOrder := squaretest.NewSquareOrder()

	order, err := square.MapOrder(squareOrder, squaretest.SquarePaymentID, squaretest.SquareMerchantID, test.TenantID, squaretest.SquareCustomerID)
	is.NoErr(err)

	expectedOrder := squaretest.NewOrder()
	err = test.Diff(expectedOrder, order)
	is.NoErr(err)
}

func TestMapOrder_FailsWithNonIntegerQuantity(t *testing.T) {
	is := is.New(t)

	squareOrder := squaretest.NewSquareOrder()
	squareOrder.OrderItems[0].Quantity = "not an int"

	_, err := square.MapOrder(squareOrder, squaretest.SquarePaymentID, squaretest.SquareMerchantID, test.TenantID, test.CustomerID)
	is.True(err != nil)
}
