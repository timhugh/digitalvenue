package square

type APIClient interface {
	GetOrder(orderId string, apiToken string) (*Order, error)
	GetCustomer(customerId string, apiToken string) (*Customer, error)
}
