package square

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventGatherer interface {
	Gather(squarePaymentID string) error
}

type eventGatherer struct {
	log          zerolog.Logger
	paymentRepo  PaymentsRepository
	merchantRepo MerchantsRepository
	orderRepo    OrdersRepository
	squareApi    Client
}

func NewEventGatherer(
	log zerolog.Logger,
	paymentRepo PaymentsRepository,
	merchantRepo MerchantsRepository,
	orderRepo OrdersRepository,
	squareApi Client,
) EventGatherer {
	return eventGatherer{
		log:          log.With().Str("caller", "eventGatherer.Gather").Logger(),
		paymentRepo:  paymentRepo,
		merchantRepo: merchantRepo,
		orderRepo:    orderRepo,
		squareApi:    squareApi,
	}
}

func (gatherer eventGatherer) Gather(squarePaymentID string) error {
	log := log.With().Str("square_payment_id", squarePaymentID).Logger()

	log.Info().Msg("Processing square payment event")

	payment, err := gatherer.paymentRepo.FindByID(squarePaymentID)
	if err != nil {
		return err
	}

	merchant, err := gatherer.merchantRepo.FindByID(payment.SquareMerchantID)
	if err != nil {
		return err
	}

	order, err := gatherer.squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)
	if err != nil {
		return err
	}

	err = gatherer.orderRepo.Create(order)
	if err != nil {
		return err
	}

	customer, err := gatherer.squareApi.GetCustomer(order.SquareCustomerID, merchant.SquareAPIToken)
	if err != nil {
		return err
	}

	log.Debug().Interface("customer", customer).Msg("Got customer details")

	// persist customer details

	// publish gathered event to SQS

	return nil
}
