package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/core"
)

func (repo *Repository) PutOrder(order core.Order) (string, error) {
	orderItems := make([]types.AttributeValue, len(order.Items))
	for i, item := range order.Items {
		var itemID string
		if item.ItemID == "" {
			itemID = repo.idGenerator.GenerateID()
		} else {
			itemID = item.ItemID
		}

		orderItems[i] = &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			ItemID: &types.AttributeValueMemberS{Value: itemID},
			Name:   &types.AttributeValueMemberS{Value: item.Name},
			Meta: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				SquareItemID: &types.AttributeValueMemberS{Value: item.Meta.SquareItemID},
			}},
		}}
	}

	var orderID string
	if order.OrderID == "" {
		orderID = repo.idGenerator.GenerateID()
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
				SquareOrderID:    &types.AttributeValueMemberS{Value: order.Meta.SquareOrderID},
				SquarePaymentID:  &types.AttributeValueMemberS{Value: order.Meta.SquarePaymentID},
				SquareMerchantID: &types.AttributeValueMemberS{Value: order.Meta.SquareMerchantID},
				SquareCustomerID: &types.AttributeValueMemberS{Value: order.Meta.SquareCustomerID},
			}},
		},
		TableName: aws.String(repo.ordersTableName),
	}

	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return "", err
	}

	return orderID, nil
}
