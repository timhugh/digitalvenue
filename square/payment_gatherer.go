package square

import (
	"context"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/logger"
)

type PaymentGatherer interface {
	Gather(ctx context.Context, payment *Payment) error
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

func (gatherer paymentGatherer) Gather(ctx context.Context, payment *Payment) error {
	_, log := logger.FromContext(ctx)
	log.Info("Processing new square payment")

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

	log.Info("Put customer: '%s'", customer.ID)

	order, err := MapOrder(squareOrder, payment.SquarePaymentID, merchant.ID, merchant.TenantID, customer.ID)
	if err != nil {
		return err
	}

	err = gatherer.orderRepo.PutOrder(order)
	if err != nil {
		return err
	}

	log.Info("Put order: '%s'", order.ID)

	return nil
}
