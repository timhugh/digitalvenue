package square_test

import (
	"github.com/go-test/deep"
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"testing"
)

func TestMapOrder(t *testing.T) {
	is := is.New(t)

	squareOrder := squaretest.NewSquareOrder()

	mapper := square.NewOrderMapper()
	order, err := mapper.MapOrder(squareOrder)
	is.NoErr(err)

	expectedOrder := squaretest.NewOrder()
	if diff := deep.Equal(order, expectedOrder); diff != nil {
		t.Error(diff)
	}
}
