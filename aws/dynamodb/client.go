package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	CustomerID = "CustomerID"
	OrderID    = "OrderID"
	TenantID   = "TenantID"
	ItemID     = "ItemID"
	Name       = "Name"
	Meta       = "Meta"

	SquareAPIToken            = "SquareAPIToken"
	SquareCustomerID          = "SquareCustomerID"
	SquareMerchantID          = "SquareMerchantID"
	SquareOrderID             = "SquareOrderID"
	SquarePaymentID           = "SquarePaymentID"
	SquareItemID              = "SquareItemID"
	SquareWebhookSignatureKey = "SquareWebhookSignatureKey"

	FirstName = "FirstName"
	LastName  = "LastName"
	Email     = "Email"
	Phone     = "Phone"

	Items = "Items"

	CustomersTableName       = "CUSTOMERS_TABLE_NAME"
	OrdersTableName          = "ORDERS_TABLE_NAME"
	SquarePaymentsTableName  = "SQUARE_PAYMENTS_TABLE_NAME"
	SquareMerchantsTableName = "SQUARE_MERCHANTS_TABLE_NAME"
)

func NewClient(config aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(config)
}
