package services

import (
	"context"
	"github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func TestTicketGenerator_GenerateTickets(t *testing.T) {
	tg := NewTicketGenerator()
	order := test.NewOrder()
	err := tg.GenerateTickets(context.TODO(), order)
	if err != nil {
		t.Errorf("Error generating tickets: %v", err)
	}
}
