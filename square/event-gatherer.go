package square

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventGatherer struct {
	log         zerolog.Logger
	paymentRepo PaymentsRepository
}

func NewEventGatherer(log zerolog.Logger, paymentRepo PaymentsRepository) EventGatherer {
	return EventGatherer{
		log:         log.With().Str("caller", "EventGatherer.Gather").Logger(),
		paymentRepo: paymentRepo,
	}
}

func (gatherer EventGatherer) Gather(squarePaymentID string) error {
	log := log.With().Str("square_payment_id", squarePaymentID).Logger()

	log.Info().Msg("Processing square payment event")

	payment, err := gatherer.paymentRepo.FindByID(squarePaymentID)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Found payment: %+v", payment)

	return nil
}
