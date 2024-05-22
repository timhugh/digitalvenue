package services

import (
	"context"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/qr"
	"strings"
	"sync"
)

type TicketGenerator struct {
	qrGenerator         core.QRCodeGenerator
	qrStore             core.QRCodeStore
	ticketRepo          core.TicketRepository
	orderProcessedQueue core.OrderProcessedQueue
}

const qrSize = 256

func NewTicketGenerator(
	qrStore core.QRCodeStore,
	ticketRepo core.TicketRepository,
	orderProcessedQueue core.OrderProcessedQueue,
) *TicketGenerator {
	return &TicketGenerator{
		qrGenerator:         qr.NewGenerator(),
		qrStore:             qrStore,
		ticketRepo:          ticketRepo,
		orderProcessedQueue: orderProcessedQueue,
	}
}

func (t *TicketGenerator) GenerateTickets(ctx context.Context, order *core.Order) error {
	tickets, err := t.buildTickets(order)
	if err != nil {
		return err
	}

	errs := make(chan error, len(tickets))
	var wg sync.WaitGroup
	for _, ticket := range tickets {
		wg.Add(1)

		go func(ticket *core.Ticket) {
			defer wg.Done()

			qrPayload := "Order#" + order.ID + " Ticket#" + ticket.ID
			qrCode, err := t.qrGenerator.Encode(qrPayload, qrSize)
			if err != nil {
				errs <- err
				return
			}
			qrCode.TenantID = order.TenantID
			qrCode.OrderID = order.ID
			qrCode.OrderItemID = ticket.ID

			url, err := t.qrStore.Save(qrCode)
			if err != nil {
				errs <- err
			}

			ticket.QRCodeURL = url
		}(ticket)
	}
	wg.Wait()

	close(errs)
	if len(errs) > 0 {
		var errorMessages []string
		for err := range errs {
			errorMessages = append(errorMessages, err.Error())
		}
		return errors.Errorf("failed to generate QR codes: %s", strings.Join(errorMessages, ", "))
	}

	err = t.ticketRepo.PutTickets(tickets)
	if err != nil {
		return errors.Wrap(err, "failed to persist tickets to repository")
	}

	err = t.orderProcessedQueue.PublishOrderProcessedEvent(order.TenantID, order.ID)
	if err != nil {
		return errors.Wrap(err, "failed to publish order processed event")
	}

	return nil
}

func (t *TicketGenerator) buildTickets(order *core.Order) ([]*core.Ticket, error) {
	tickets := make([]*core.Ticket, len(order.Items))
	for i, item := range order.Items {
		tickets[i] = &core.Ticket{
			ID:         item.ID,
			TenantID:   order.TenantID,
			OrderID:    order.ID,
			CustomerID: order.CustomerID,
			Name:       item.Name,
		}
	}
	return tickets, nil
}
