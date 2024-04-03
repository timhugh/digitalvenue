package square

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
)

type EventGatherer interface {
	Gather(squarePaymentID string) error
}

type eventGatherer struct {
	log          zerolog.Logger
	paymentRepo  PaymentsRepository
	merchantRepo MerchantsRepository
	orderRepo    OrdersRepository
	customerRepo core.CustomerRepository
	squareApi    Client
}

func NewEventGatherer(
	log zerolog.Logger,
	paymentRepo PaymentsRepository,
	merchantRepo MerchantsRepository,
	orderRepo OrdersRepository,
	customerRepo core.CustomerRepository,
	squareApi Client,
) EventGatherer {
	return eventGatherer{
		log:          log.With().Str("caller", "eventGatherer.Gather").Logger(),
		paymentRepo:  paymentRepo,
		merchantRepo: merchantRepo,
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
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

	squareCustomer, err := gatherer.squareApi.GetCustomer(order.SquareCustomerID, merchant.SquareAPIToken)
	if err != nil {
		return err
	}

	customer := MapCustomer(squareCustomer)

	err = gatherer.customerRepo.Create(customer)
	if err != nil {
		return err
	}

	// publish gathered event to SQS

	return nil
}
