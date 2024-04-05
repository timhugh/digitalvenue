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
	paymentRepo  PaymentRepository
	merchantRepo MerchantRepository
	orderRepo    core.OrderRepository
	customerRepo core.CustomerRepository
	squareApi    Client
	orderMapper  OrderMapper
}

func NewEventGatherer(
	log zerolog.Logger,
	paymentRepo PaymentRepository,
	merchantRepo MerchantRepository,
	orderRepo core.OrderRepository,
	customerRepo core.CustomerRepository,
	squareApi Client,
	orderMapper OrderMapper,
) EventGatherer {
	return eventGatherer{
		log:          log.With().Str("caller", "eventGatherer.Gather").Logger(),
		paymentRepo:  paymentRepo,
		merchantRepo: merchantRepo,
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		squareApi:    squareApi,
		orderMapper:  orderMapper,
	}
}

func (gatherer eventGatherer) Gather(squarePaymentID string) error {
	log := log.With().Str("square_payment_id", squarePaymentID).Logger()

	log.Info().Msg("Processing square payment event")

	payment, err := gatherer.paymentRepo.GetSquarePayment(squarePaymentID)
	if err != nil {
		return err
	}

	merchant, err := gatherer.merchantRepo.GetSquareMerchant(payment.SquareMerchantID)
	if err != nil {
		return err
	}

	squareOrder, err := gatherer.squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)
	if err != nil {
		return err
	}

	squareCustomer, err := gatherer.squareApi.GetCustomer(squareOrder.SquareCustomerID, merchant.SquareAPIToken)
	if err != nil {
		return err
	}

	customer := MapCustomer(squareCustomer)
	customerID, err := gatherer.customerRepo.PutCustomer(customer)
	if err != nil {
		return err
	}

	order, err := gatherer.orderMapper.MapOrder(squareOrder)
	if err != nil {
		return err
	}
	order.CustomerID = customerID
	order.Meta.SquarePaymentID = squarePaymentID
	order.Meta.SquareMerchantID = merchant.SquareMerchantID
	orderID, err := gatherer.orderRepo.PutOrder(order)
	if err != nil {
		return err
	}

	log.Info().Str("order_id", orderID).Msg("Order created")

	// publish gathered event to SQS

	return nil
}
