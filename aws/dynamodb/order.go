package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/aws/dynamodb/square"
	"github.com/timhugh/digitalvenue/core"
	"os"
)

type OrderRepositoryConfig struct {
	TableName string
}

func NewOrderRepositoryConfig() OrderRepositoryConfig {
	return OrderRepositoryConfig{
		TableName: os.Getenv(OrdersTableName),
	}
}

type OrderRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewOrderRepository(config OrderRepositoryConfig, client *dynamodb.Client) core.OrderRepository {
	return OrderRepository{
		tableName: config.TableName,
		client:    client,
	}
}

func (repo OrderRepository) Put(order core.Order) (string, error) {

	orderItems := make([]types.AttributeValue, len(order.Items))
	for i, item := range order.Items {
		var itemID string
		if item.ItemID == "" {
			itemID = core.GenerateID()
		} else {
			itemID = item.ItemID
		}

		orderItems[i] = &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			ItemID: &types.AttributeValueMemberS{Value: itemID},
			Name:   &types.AttributeValueMemberS{Value: item.Name},
			Meta: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				square.SquareItemID: &types.AttributeValueMemberS{Value: item.Meta.SquareItemID},
			}},
		}}
	}

	var orderID string
	if order.OrderID == "" {
		orderID = core.GenerateID()
	} else {
		orderID = order.OrderID
	}

	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			OrderID:    &types.AttributeValueMemberS{Value: orderID},
			TenantID:   &types.AttributeValueMemberS{Value: order.TenantID},
			CustomerID: &types.AttributeValueMemberS{Value: order.CustomerID},
			Items:      &types.AttributeValueMemberL{Value: orderItems},
			Meta: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				square.SquareOrderID:    &types.AttributeValueMemberS{Value: order.Meta.SquareOrderID},
				square.SquarePaymentID:  &types.AttributeValueMemberS{Value: order.Meta.SquarePaymentID},
				square.SquareMerchantID: &types.AttributeValueMemberS{Value: order.Meta.SquareMerchantID},
				square.SquareCustomerID: &types.AttributeValueMemberS{Value: order.Meta.SquareCustomerID},
			}},
		},
		TableName: aws.String(repo.tableName),
	}

	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return "", err
	}

	return orderID, nil
}
