package dynamodb

import (
	"context"
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

	CustomersTableNameKey       = "CUSTOMERS_TABLE_NAME"
	OrdersTableNameKey          = "ORDERS_TABLE_NAME"
	SquarePaymentsTableNameKey  = "SQUARE_PAYMENTS_TABLE_NAME"
	SquareMerchantsTableNameKey = "SQUARE_MERCHANTS_TABLE_NAME"
)

type Client interface {
	PutItem(ctx context.Context, input *dynamodb.PutItemInput, optFns ...func(options *dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, input *dynamodb.GetItemInput, optFns ...func(options *dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func NewClient(config aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(config)
}
