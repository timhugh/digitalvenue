package square

import (
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core"
)

type PaymentGatherer interface {
	Gather(payment *Payment, log zerolog.Logger) error
}

type paymentGatherer struct {
	merchantRepo MerchantRepository
	orderRepo    core.OrderRepository
	customerRepo core.CustomerRepository
	squareApi    APIClient
}

func NewPaymentGatherer(
	merchantRepo MerchantRepository,
	orderRepo core.OrderRepository,
	customerRepo core.CustomerRepository,
	squareApi APIClient,
) PaymentGatherer {
	return paymentGatherer{
		merchantRepo: merchantRepo,
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		squareApi:    squareApi,
	}
}

func (gatherer paymentGatherer) Gather(payment *Payment, log zerolog.Logger) error {
	log.Info().Msg("Processing new square payment")

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

	customer := MapCustomer(squareCustomer, merchant.TenantID)
	err = gatherer.customerRepo.PutCustomer(customer)
	if err != nil {
		return err
	}

	order, err := MapOrder(squareOrder, payment.SquarePaymentID, merchant.ID, merchant.TenantID, customer.ID)
	if err != nil {
		return err
	}

	err = gatherer.orderRepo.PutOrder(order)
	if err != nil {
		return err
	}

	log.Info().Str("order_id", order.ID).Msg("Order created")

	return nil
}
