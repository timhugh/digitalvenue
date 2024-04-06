package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/core"
)

type OrderRepository struct {
	client      Client
	idGenerator core.IDGenerator
	tableName   string
}

func NewOrderRepository(client Client) (*OrderRepository, error) {
	tableName, err := core.RequireEnv(OrdersTableNameKey)
	if err != nil {
		return nil, err
	}

	return &OrderRepository{
		client:      client,
		idGenerator: core.NewIDGenerator(),
		tableName:   tableName,
	}, nil
}

func (repo *OrderRepository) PutOrder(order core.Order) (string, error) {
	orderItems := make([]types.AttributeValue, len(order.Items))
	for i, item := range order.Items {
		var itemID string
		if item.ItemID == "" {
			itemID = repo.idGenerator.GenerateID()
		} else {
			itemID = item.ItemID
		}

		meta := make(map[string]types.AttributeValue)
		for key, value := range item.Meta {
			meta[key] = &types.AttributeValueMemberS{Value: value}
		}

		orderItems[i] = &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			ItemID: &types.AttributeValueMemberS{Value: itemID},
			Name:   &types.AttributeValueMemberS{Value: item.Name},
			Meta:   &types.AttributeValueMemberM{Value: meta},
		}}
	}

	var orderID string
	if order.OrderID == "" {
		orderID = repo.idGenerator.GenerateID()
	} else {
		orderID = order.OrderID
	}

	meta := make(map[string]types.AttributeValue)
	for key, value := range order.Meta {
		meta[key] = &types.AttributeValueMemberS{Value: value}
	}

	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			OrderID:    &types.AttributeValueMemberS{Value: orderID},
			TenantID:   &types.AttributeValueMemberS{Value: order.TenantID},
			CustomerID: &types.AttributeValueMemberS{Value: order.CustomerID},
			Items:      &types.AttributeValueMemberL{Value: orderItems},
			Meta:       &types.AttributeValueMemberM{Value: meta},
		},
		TableName: aws.String(repo.tableName),
	}

	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return "", err
	}

	return orderID, nil
}
