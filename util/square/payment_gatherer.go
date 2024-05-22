package square

import (
	"context"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
)

type PaymentGatherer interface {
	Gather(ctx context.Context, squareMerchantID string, squarePaymentID string) error
}

type paymentGatherer struct {
	paymentRepo       PaymentRepository
	merchantRepo      MerchantRepository
	orderRepo         core.OrderRepository
	customerRepo      core.CustomerRepository
	squareApi         APIClient
	orderCreatedQueue core.OrderCreatedQueue
	orderBuilder      *OrderBuilder
}

func NewPaymentGatherer(
	paymentRepo PaymentRepository,
	merchantRepo MerchantRepository,
	orderRepo core.OrderRepository,
	customerRepo core.CustomerRepository,
	squareApi APIClient,
	orderCreatedQueue core.OrderCreatedQueue,
	orderBuilder *OrderBuilder,
) PaymentGatherer {
	return paymentGatherer{
		paymentRepo:       paymentRepo,
		merchantRepo:      merchantRepo,
		orderRepo:         orderRepo,
		customerRepo:      customerRepo,
		squareApi:         squareApi,
		orderCreatedQueue: orderCreatedQueue,
		orderBuilder:      orderBuilder,
	}
}

func (gatherer paymentGatherer) Gather(ctx context.Context, squareMerchantID string, squarePaymentID string) error {
	_, log := logger.FromContext(ctx)
	log = log.Sub().AddParams(map[string]interface{}{
		"squareMerchantId": squareMerchantID,
		"squarePaymentId":  squarePaymentID,
	})
	log.Info("Processing new square payment")

	payment, err := gatherer.paymentRepo.GetSquarePayment(squareMerchantID, squarePaymentID)
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

	log.AddParam("squareCustomerID", squareOrder.SquareCustomerID)

	var customer *core.Customer
	existingCustomer, err := gatherer.customerRepo.GetCustomer(merchant.TenantID, squareOrder.SquareCustomerID)
	if err != nil {
		_, ok := err.(core.ItemNotFoundException)
		if !ok {
			return err
		}
		// TODO: this is a little gross, just a quick fix
		existingCustomer = nil
	}

	if existingCustomer == nil {
		log.Info("Creating new customer for square payment")
		squareCustomer, err := gatherer.squareApi.GetCustomer(squareOrder.SquareCustomerID, merchant.SquareAPIToken)
		if err != nil {
			return err
		}

		customer = MapCustomer(squareCustomer, merchant.TenantID)
		err = gatherer.customerRepo.PutCustomer(customer)
		if err != nil {
			return err
		}
	} else {
		log.Info("Found existing customer for square payment")
		customer = existingCustomer
	}

	existingOrder, err := gatherer.orderRepo.GetOrder(merchant.TenantID, squareOrder.SquareOrderID)
	if err != nil {
		_, ok := err.(core.ItemNotFoundException)
		if !ok {
			return err
		}
		// TODO: this is a little gross, just a quick fix
		existingOrder = nil
	}

	if existingOrder != nil {
		log.Info("Found existing order for square payment; skipping processing")
		return nil
	}

	log.Info("Creating new order for square payment")

	order, err := gatherer.orderBuilder.BuildOrder(squareOrder, merchant, payment.SquarePaymentID, customer.ID)
	if err != nil {
		return err
	}

	err = gatherer.orderRepo.PutOrder(order)
	if err != nil {
		return err
	}

	log.Info("Publishing order created event for square payment")

	err = gatherer.orderCreatedQueue.PublishOrderCreatedEvent(order.TenantID, order.ID)
	if err != nil {
		return err
	}

	return nil
}
