package services

import (
	"context"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func TestTicketGenerator_GenerateTickets(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	qrStorer := mock.Mock[core.QRCodeStorer]()
	ticketRepo := mock.Mock[core.TicketRepository]()

	tg := NewTicketGenerator(qrStorer, ticketRepo)
	order := test.NewOrder()

	err := tg.GenerateTickets(context.TODO(), order)
	is.NoErr(err)
}
