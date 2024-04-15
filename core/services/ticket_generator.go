package services

import (
	"context"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/logger"
	"github.com/timhugh/digitalvenue/qr"
)

type TicketGenerator struct {
	qrGenerator core.QRCodeGenerator
}

const qrSize = 256

func NewTicketGenerator() *TicketGenerator {
	return &TicketGenerator{
		qrGenerator: qr.NewGenerator(),
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

	for _, ticket := range tickets {
		qrPayload := "Order#" + order.ID + " Ticket#" + ticket.ID
		qrCode, err := t.qrGenerator.Encode(qrPayload, qrSize)
		if err != nil {
			return err
		}
		ticket.QRCodeBase64 = qrCode.Image
		ticket.QRCodeFileType = qrCode.FileType

		log.Sub().AddParam("qr_code", ticket.QRCodeBase64).Debug("Generated qrcode for ticket %s", ticket.ID)
	}

	log.Info("Finished building and generating tickets")

	return nil
}

func buildTickets(ctx context.Context, order *core.Order) ([]*core.Ticket, error) {
	_, log := logger.FromContext(ctx)

	tickets := make([]*core.Ticket, len(order.Items))
	for i, item := range order.Items {
		ticket := buildTicket(order, &item)
		tickets[i] = ticket
		log.Debug("Built ticket %s from order item %s", ticket.ID, item.ID)
	}
	return tickets, nil
}

func buildTicket(order *core.Order, item *core.OrderItem) *core.Ticket {
	return &core.Ticket{
		ID:         item.ID,
		TenantID:   order.TenantID,
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Name:       item.Name,
	}
}
