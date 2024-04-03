package square

import (
	"fmt"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core"
	"testing"
)

func TestEventGatherer_GatherSuccess(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	log := zerolog.Logger{}
	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	expectedOrder := Order{
		SquareOrderID:    "squareOrderID",
		SquareCustomerID: "squareCustomerID",
	}
	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)

	orderRepo := mock.Mock[OrdersRepository]()
	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(nil)

	expectedSquareCustomer := Customer{
		SquareCustomerID: "squareCustomerID",
		FirstName:        "Tim",
		LastName:         "Heuett",
		Email:            "info@timheuett.com",
		Phone:            "+12062062062",
	}
	mock.WhenDouble(squareApi.GetCustomer(expectedOrder.SquareCustomerID, merchant.SquareAPIToken)).ThenReturn(expectedSquareCustomer, nil)

	expectedCustomer := core.Customer{
		CustomerID: "",
		FirstName:  "Tim",
		LastName:   "Heuett",
		Email:      "info@timheuett.com",
		Phone:      "+12062062062",
		Meta: core.CustomerMeta{
			Source:           "square",
			SquareCustomerID: "squareCustomerID",
		},
	}
	customerRepo := mock.Mock[core.CustomerRepository]()
	mock.WhenSingle(customerRepo.Create(mock.Equal(expectedCustomer))).ThenReturn(nil)

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)

	err := gatherer.Gather("squarePaymentID")
	is.NoErr(err)
}

func TestEventGatherer_GatherPaymentRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	squarePaymentID := "squarePaymentID"

	log := zerolog.Logger{}

	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(squarePaymentID)).ThenReturn(Payment{}, fmt.Errorf("failed to find payment"))

	merchantRepo := mock.Mock[MerchantsRepository]()
	squareApi := mock.Mock[Client]()
	orderRepo := mock.Mock[OrdersRepository]()
	customerRepo := mock.Mock[core.CustomerRepository]()

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)

	err := gatherer.Gather(squarePaymentID)
	is.Equal(err.Error(), "failed to find payment")
}

func TestEventGatherer_GatherMerchantRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	squarePaymentID := "squarePaymentID"

	log := zerolog.Logger{}

	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(Merchant{}, fmt.Errorf("failed to find merchant"))

	squareApi := mock.Mock[Client]()
	orderRepo := mock.Mock[OrdersRepository]()
	customerRepo := mock.Mock[core.CustomerRepository]()

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)

	err := gatherer.Gather(squarePaymentID)
	is.Equal(err.Error(), "failed to find merchant")
}

func TestEventGatherer_GatherSquareFetchOrderError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	squarePaymentID := "squarePaymentID"

	log := zerolog.Logger{}

	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(Order{}, fmt.Errorf("failed to fetch order"))

	orderRepo := mock.Mock[OrdersRepository]()
	customerRepo := mock.Mock[core.CustomerRepository]()

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)

	err := gatherer.Gather(squarePaymentID)
	is.Equal(err.Error(), "failed to fetch order")
}

func TestEventGatherer_GatherOrderRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	log := zerolog.Logger{}
	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	expectedOrder := Order{
		SquareOrderID:    "squareOrderID",
		SquareCustomerID: "squareCustomerID",
	}
	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)
	orderRepo := mock.Mock[OrdersRepository]()
	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(fmt.Errorf("failed to create order"))

	customerRepo := mock.Mock[core.CustomerRepository]()

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)

	err := gatherer.Gather("squarePaymentID")
	is.Equal(err.Error(), "failed to create order")
}

func TestEventGatherer_GatherSquareFetchCustomerError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	log := zerolog.Logger{}
	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	expectedOrder := Order{
		SquareOrderID:    "squareOrderID",
		SquareCustomerID: "squareCustomerID",
	}
	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)

	orderRepo := mock.Mock[OrdersRepository]()
	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(nil)

	mock.WhenDouble(squareApi.GetCustomer(expectedOrder.SquareCustomerID, merchant.SquareAPIToken)).ThenReturn(Customer{}, fmt.Errorf("failed to fetch customer"))

	customerRepo := mock.Mock[core.CustomerRepository]()

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)

	err := gatherer.Gather("squarePaymentID")
	is.Equal(err.Error(), "failed to fetch customer")
}

func TestEventGatherer_GatherCustomerRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	log := zerolog.Logger{}
	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	expectedOrder := Order{
		SquareOrderID:    "squareOrderID",
		SquareCustomerID: "squareCustomerID",
	}
	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)

	orderRepo := mock.Mock[OrdersRepository]()
	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(nil)

	expectedSquareCustomer := Customer{
		SquareCustomerID: "squareCustomerID",
		FirstName:        "Tim",
		LastName:         "Heuett",
		Email:            "info@timheuett.com",
		Phone:            "+12062062062",
	}
	mock.WhenDouble(squareApi.GetCustomer(expectedOrder.SquareCustomerID, merchant.SquareAPIToken)).ThenReturn(expectedSquareCustomer, nil)

	expectedCustomer := core.Customer{
		CustomerID: "",
		FirstName:  "Tim",
		LastName:   "Heuett",
		Email:      "info@timheuett.com",
		Phone:      "+12062062062",
		Meta: core.CustomerMeta{
			Source:           "square",
			SquareCustomerID: "squareCustomerID",
		},
	}
	customerRepo := mock.Mock[core.CustomerRepository]()
	mock.WhenSingle(customerRepo.Create(mock.Equal(expectedCustomer))).ThenReturn(fmt.Errorf("failed to create customer"))

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)

	err := gatherer.Gather("squarePaymentID")
	is.Equal(err.Error(), "failed to create customer")
}
