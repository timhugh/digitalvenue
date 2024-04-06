package square_test

import (
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"testing"
)

func TestMapCustomer(t *testing.T) {
	is := is.New(t)

	squareCustomer := squaretest.NewSquareCustomer()

	customer := square.MapCustomer(squareCustomer)

	expectedCustomer := squaretest.NewCustomer()
	is.Equal(customer, expectedCustomer)
}
