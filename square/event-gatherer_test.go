package square

import (
	"github.com/go-test/deep"
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
	mock.WhenDouble(paymentRepo.Get(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.Get(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	squareOrder := Order{
		SquareOrderID:    "squareOrderID",
		SquareCustomerID: "squareCustomerID",
		OrderItems: []OrderItem{
			{
				ItemID:   "squareItemID",
				Name:     "item1",
				Quantity: "1",
			},
		},
	}
	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(squareOrder, nil)

	expectedOrder := core.Order{
		CustomerID: "customer_id",
		Items: []core.OrderItem{
			{
				Name: "item1",
				Meta: core.OrderItemMeta{
					SquareItemID: "squareItemID",
				},
			},
		},
		Meta: core.OrderMeta{
			SquareOrderID:    "squareOrderID",
			SquareCustomerID: "squareCustomerID",
			SquarePaymentID:  "squarePaymentID",
			SquareMerchantID: "squareMerchantID",
		},
	}
	orderRepo := mock.Mock[core.OrderRepository]()
	orderCaptor := mock.Captor[core.Order]()
	mock.WhenDouble(orderRepo.Put(orderCaptor.Capture())).ThenReturn("order_id", nil)

	expectedSquareCustomer := Customer{
		SquareCustomerID: "squareCustomerID",
		FirstName:        "Tim",
		LastName:         "Heuett",
		Email:            "info@timheuett.com",
		Phone:            "+12062062062",
	}
	mock.WhenDouble(squareApi.GetCustomer(squareOrder.SquareCustomerID, merchant.SquareAPIToken)).ThenReturn(expectedSquareCustomer, nil)

	expectedCustomer := core.Customer{
		FirstName: "Tim",
		LastName:  "Heuett",
		Email:     "info@timheuett.com",
		Phone:     "+12062062062",
		Meta: core.CustomerMeta{
			Source:           "square",
			SquareCustomerID: "squareCustomerID",
		},
	}
	customerRepo := mock.Mock[core.CustomerRepository]()
	customerCaptor := mock.Captor[core.Customer]()
	mock.WhenDouble(customerRepo.Put(customerCaptor.Capture())).ThenReturn(expectedOrder.CustomerID, nil)

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)
	err := gatherer.Gather(payment.SquarePaymentID)
	is.NoErr(err)

	persistedOrder := orderCaptor.Last()
	if diff := deep.Equal(expectedOrder, persistedOrder); diff != nil {
		t.Error(diff)
	}

	persistedCustomer := customerCaptor.Last()
	if diff := deep.Equal(expectedCustomer, persistedCustomer); diff != nil {
		t.Error(diff)
	}
}

//
//func TestEventGatherer_GatherPaymentRepoError(t *testing.T) {
//	is := is.New(t)
//	mock.SetUp(t)
//
//	squarePaymentID := "squarePaymentID"
//
//	log := zerolog.Logger{}
//
//	paymentRepo := mock.Mock[PaymentsRepository]()
//	mock.WhenDouble(paymentRepo.Get(squarePaymentID)).ThenReturn(Payment{}, fmt.Errorf("failed to find payment"))
//
//	merchantRepo := mock.Mock[MerchantsRepository]()
//	squareApi := mock.Mock[Client]()
//	orderRepo := mock.Mock[core.OrderRepository]()
//	customerRepo := mock.Mock[core.CustomerRepository]()
//
//	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)
//
//	err := gatherer.Gather(squarePaymentID)
//	is.Equal(err.Error(), "failed to find payment")
//}
//
//func TestEventGatherer_GatherMerchantRepoError(t *testing.T) {
//	is := is.New(t)
//	mock.SetUp(t)
//
//	squarePaymentID := "squarePaymentID"
//
//	log := zerolog.Logger{}
//
//	payment := Payment{
//		SquareMerchantID: "squareMerchantID",
//		SquareOrderID:    "squareOrderID",
//		SquarePaymentID:  "squarePaymentID",
//	}
//	paymentRepo := mock.Mock[PaymentsRepository]()
//	mock.WhenDouble(paymentRepo.Get(payment.SquarePaymentID)).ThenReturn(payment, nil)
//
//	merchantRepo := mock.Mock[MerchantsRepository]()
//	mock.WhenDouble(merchantRepo.Get(payment.SquareMerchantID)).ThenReturn(Merchant{}, fmt.Errorf("failed to find merchant"))
//
//	squareApi := mock.Mock[Client]()
//	orderRepo := mock.Mock[core.OrderRepository]()
//	customerRepo := mock.Mock[core.CustomerRepository]()
//
//	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)
//
//	err := gatherer.Gather(squarePaymentID)
//	is.Equal(err.Error(), "failed to find merchant")
//}
//
//func TestEventGatherer_GatherSquareFetchOrderError(t *testing.T) {
//	is := is.New(t)
//	mock.SetUp(t)
//
//	squarePaymentID := "squarePaymentID"
//
//	log := zerolog.Logger{}
//
//	payment := Payment{
//		SquareMerchantID: "squareMerchantID",
//		SquareOrderID:    "squareOrderID",
//		SquarePaymentID:  "squarePaymentID",
//	}
//	paymentRepo := mock.Mock[PaymentsRepository]()
//	mock.WhenDouble(paymentRepo.Get(payment.SquarePaymentID)).ThenReturn(payment, nil)
//
//	merchant := Merchant{
//		SquareMerchantID: "squareMerchantID",
//		SquareAPIToken:   "squareAPIToken",
//	}
//	merchantRepo := mock.Mock[MerchantsRepository]()
//	mock.WhenDouble(merchantRepo.Get(payment.SquareMerchantID)).ThenReturn(merchant, nil)
//
//	squareApi := mock.Mock[Client]()
//	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(Order{}, fmt.Errorf("failed to fetch order"))
//
//	orderRepo := mock.Mock[OrderRepository]()
//	customerRepo := mock.Mock[core.CustomerRepository]()
//
//	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)
//
//	err := gatherer.Gather(squarePaymentID)
//	is.Equal(err.Error(), "failed to fetch order")
//}
//
//func TestEventGatherer_GatherOrderRepoError(t *testing.T) {
//	is := is.New(t)
//	mock.SetUp(t)
//
//	log := zerolog.Logger{}
//	payment := Payment{
//		SquareMerchantID: "squareMerchantID",
//		SquareOrderID:    "squareOrderID",
//		SquarePaymentID:  "squarePaymentID",
//	}
//	paymentRepo := mock.Mock[PaymentsRepository]()
//	mock.WhenDouble(paymentRepo.Get(payment.SquarePaymentID)).ThenReturn(payment, nil)
//
//	merchant := Merchant{
//		SquareMerchantID: "squareMerchantID",
//		SquareAPIToken:   "squareAPIToken",
//	}
//	merchantRepo := mock.Mock[MerchantsRepository]()
//	mock.WhenDouble(merchantRepo.Get(payment.SquareMerchantID)).ThenReturn(merchant, nil)
//
//	expectedOrder := Order{
//		SquareOrderID:    "squareOrderID",
//		SquareCustomerID: "squareCustomerID",
//	}
//	squareApi := mock.Mock[Client]()
//	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)
//	orderRepo := mock.Mock[OrderRepository]()
//	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(fmt.Errorf("failed to create order"))
//
//	customerRepo := mock.Mock[core.CustomerRepository]()
//
//	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)
//
//	err := gatherer.Gather("squarePaymentID")
//	is.Equal(err.Error(), "failed to create order")
//}
//
//func TestEventGatherer_GatherSquareFetchCustomerError(t *testing.T) {
//	is := is.New(t)
//	mock.SetUp(t)
//
//	log := zerolog.Logger{}
//	payment := Payment{
//		SquareMerchantID: "squareMerchantID",
//		SquareOrderID:    "squareOrderID",
//		SquarePaymentID:  "squarePaymentID",
//	}
//	paymentRepo := mock.Mock[PaymentsRepository]()
//	mock.WhenDouble(paymentRepo.Get(payment.SquarePaymentID)).ThenReturn(payment, nil)
//
//	merchant := Merchant{
//		SquareMerchantID: "squareMerchantID",
//		SquareAPIToken:   "squareAPIToken",
//	}
//	merchantRepo := mock.Mock[MerchantsRepository]()
//	mock.WhenDouble(merchantRepo.Get(payment.SquareMerchantID)).ThenReturn(merchant, nil)
//
//	expectedOrder := Order{
//		SquareOrderID:    "squareOrderID",
//		SquareCustomerID: "squareCustomerID",
//	}
//	squareApi := mock.Mock[Client]()
//	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)
//
//	orderRepo := mock.Mock[OrderRepository]()
//	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(nil)
//
//	mock.WhenDouble(squareApi.GetCustomer(expectedOrder.SquareCustomerID, merchant.SquareAPIToken)).ThenReturn(Customer{}, fmt.Errorf("failed to fetch customer"))
//
//	customerRepo := mock.Mock[core.CustomerRepository]()
//
//	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)
//
//	err := gatherer.Gather("squarePaymentID")
//	is.Equal(err.Error(), "failed to fetch customer")
//}
//
//func TestEventGatherer_GatherCustomerRepoError(t *testing.T) {
//	is := is.New(t)
//	mock.SetUp(t)
//
//	log := zerolog.Logger{}
//	payment := Payment{
//		SquareMerchantID: "squareMerchantID",
//		SquareOrderID:    "squareOrderID",
//		SquarePaymentID:  "squarePaymentID",
//	}
//	paymentRepo := mock.Mock[PaymentsRepository]()
//	mock.WhenDouble(paymentRepo.Get(payment.SquarePaymentID)).ThenReturn(payment, nil)
//
//	merchant := Merchant{
//		SquareMerchantID: "squareMerchantID",
//		SquareAPIToken:   "squareAPIToken",
//	}
//	merchantRepo := mock.Mock[MerchantsRepository]()
//	mock.WhenDouble(merchantRepo.Get(payment.SquareMerchantID)).ThenReturn(merchant, nil)
//
//	expectedOrder := Order{
//		SquareOrderID:    "squareOrderID",
//		SquareCustomerID: "squareCustomerID",
//	}
//	squareApi := mock.Mock[Client]()
//	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)
//
//	orderRepo := mock.Mock[OrderRepository]()
//	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(nil)
//
//	expectedSquareCustomer := Customer{
//		SquareCustomerID: "squareCustomerID",
//		FirstName:        "Tim",
//		LastName:         "Heuett",
//		Email:            "info@timheuett.com",
//		Phone:            "+12062062062",
//	}
//	mock.WhenDouble(squareApi.GetCustomer(expectedOrder.SquareCustomerID, merchant.SquareAPIToken)).ThenReturn(expectedSquareCustomer, nil)
//
//	expectedCustomer := core.Customer{
//		CustomerID: "",
//		FirstName:  "Tim",
//		LastName:   "Heuett",
//		Email:      "info@timheuett.com",
//		Phone:      "+12062062062",
//		Meta: core.CustomerMeta{
//			Source:           "square",
//			SquareCustomerID: "squareCustomerID",
//		},
//	}
//	customerRepo := mock.Mock[core.CustomerRepository]()
//	mock.WhenSingle(customerRepo.Put(mock.Equal(expectedCustomer))).ThenReturn(fmt.Errorf("failed to create customer"))
//
//	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)
//
//	err := gatherer.Gather("squarePaymentID")
//	is.Equal(err.Error(), "failed to create customer")
//}
