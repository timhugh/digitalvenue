package square_test

import (
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/util/square"
	"github.com/timhugh/digitalvenue/util/square/squaretest"
	"github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func TestMapCustomer(t *testing.T) {
	is := is.New(t)

	squareCustomer := squaretest.NewSquareCustomer()

	customer := square.MapCustomer(squareCustomer, test.TenantID)

	expectedCustomer := squaretest.NewCustomer()
	is.Equal(customer, expectedCustomer)
}
