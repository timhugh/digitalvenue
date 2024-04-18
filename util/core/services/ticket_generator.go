package services

import (
	"context"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/qr"
	"strings"
	"sync"
)

type TicketGenerator struct {
	qrGenerator core.QRCodeGenerator
	qrStore     core.QRCodeStorer
	ticketRepo  core.TicketRepository
}

const qrSize = 256

func NewTicketGenerator(qrStore core.QRCodeStorer, ticketRepo core.TicketRepository) *TicketGenerator {
	return &TicketGenerator{
		qrGenerator: qr.NewGenerator(),
		qrStore:     qrStore,
		ticketRepo:  ticketRepo,
	}
}

func (t *TicketGenerator) GenerateTickets(ctx context.Context, order *core.Order) error {
	ctx, log := logger.FromContext(ctx)

	log.Info("Building tickets from order")

	tickets, err := buildTickets(ctx, order)
	if err != nil {
		return err
	}

	log.Debug("Generating QR codes for tickets")

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

			log.Sub().AddParam("qr_code", qrCode.Image).Debug("Generated qrcode for ticket %s", ticket.ID)

			url, err := t.qrStore.Save(qrCode)
			if err != nil {
				errs <- err
			}

			ticket.QRCodeURL = url

			log.Debug("QR code for ticket %s saved to %s", ticket.ID, url)
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

	log.Debug("Successfully generated and saved QR codes for tickets")

	err = t.ticketRepo.PutTickets(tickets)
	if err != nil {
		return errors.Wrap(err, "failed to persist tickets to repository")
	}

	return nil
}

func buildTickets(ctx context.Context, order *core.Order) ([]*core.Ticket, error) {
	_, log := logger.FromContext(ctx)

	tickets := make([]*core.Ticket, len(order.Items))
	for i, item := range order.Items {
		tickets[i] = &core.Ticket{
			ID:         item.ID,
			TenantID:   order.TenantID,
			OrderID:    order.ID,
			CustomerID: order.CustomerID,
			Name:       item.Name,
		}
		log.Debug("Built ticket %s", tickets[i].ID)
	}
	return tickets, nil
}
