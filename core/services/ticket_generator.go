package services

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core"
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
	tickets, err := buildTickets(ctx, order)
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		qrPayload := "Order#" + order.ID + " Ticket#" + ticket.ID
		qrCode, err := t.qrGenerator.Encode(qrPayload, qrSize)
		if err != nil {
			return err
		}
		ticket.QRCodeBase64 = qrCode.Image
		ticket.QRCodeFileType = qrCode.FileType
	}

	return nil
}

func buildTickets(ctx context.Context, order *core.Order) ([]*core.Ticket, error) {
	logger := zerolog.Ctx(ctx)

	tickets := make([]*core.Ticket, len(order.Items))
	for i, item := range order.Items {
		ticket := buildTicket(order, &item)
		tickets[i] = ticket
		logger.Debug().Msgf("Built ticket %s for order %s", ticket.ID, order.ID)
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
