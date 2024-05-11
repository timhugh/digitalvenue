package square

import (
	"context"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
)

type PaymentGatherer interface {
	Gather(ctx context.Context, payment *Payment) error
}

type paymentGatherer struct {
	merchantRepo      MerchantRepository
	orderRepo         core.OrderRepository
	customerRepo      core.CustomerRepository
	squareApi         APIClient
	orderCreatedQueue core.OrderCreatedQueue
}

func NewPaymentGatherer(
	merchantRepo MerchantRepository,
	orderRepo core.OrderRepository,
	customerRepo core.CustomerRepository,
	squareApi APIClient,
	orderCreatedQueue core.OrderCreatedQueue,
) PaymentGatherer {
	return paymentGatherer{
		merchantRepo:      merchantRepo,
		orderRepo:         orderRepo,
		customerRepo:      customerRepo,
		squareApi:         squareApi,
		orderCreatedQueue: orderCreatedQueue,
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

	log.Debug("Put customer: '%s'", customer.ID)

	order, err := MapOrder(squareOrder, payment.SquarePaymentID, merchant.ID, merchant.TenantID, customer.ID)
	if err != nil {
		return err
	}

	err = gatherer.orderRepo.PutOrder(order)
	if err != nil {
		return err
	}

	log.Debug("Put order: '%s'", order.ID)

	err = gatherer.orderCreatedQueue.PublishOrderCreatedEvent(order.TenantID, order.ID)
	if err != nil {
		return err
	}

	return nil
}
