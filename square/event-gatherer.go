package square

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/square/db"
)

type SquareEventGatherer struct {
	log         zerolog.Logger
	paymentRepo db.SquarePaymentsRepository
}

func (gatherer SquareEventGatherer) Gather(squarePaymentID string) error {
	log := log.With().
		Str("payment_id", squarePaymentID).
		Str("caller", "SquareEventGatherer.Gather").
		Logger()

	log.Info().Msgf("Processing payment event")

	payment, err := gatherer.paymentRepo.FindByID(squarePaymentID)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Found payment: %+v", payment)

	return nil
}
