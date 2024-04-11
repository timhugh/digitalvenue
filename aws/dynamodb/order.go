package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/core"
)

func (repo *Repository) PutOrder(order core.Order) error {
	tenantKey := "Tenant#" + order.TenantID
	orderKey := "Order#" + order.ID

	orderItems := make([]types.AttributeValue, len(order.Items))
	for i, item := range order.Items {
		var meta map[string]types.AttributeValue
		if item.Meta != nil {
			meta = make(map[string]types.AttributeValue)
			for k, v := range item.Meta {
				meta[k] = &types.AttributeValueMemberS{Value: v}
			}
		}

		orderItems[i] = &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"ItemID": &types.AttributeValueMemberS{Value: item.ID},
				"Name":   &types.AttributeValueMemberS{Value: item.Name},
				"Meta":   &types.AttributeValueMemberM{Value: meta},
			},
		}
	}

	var meta map[string]types.AttributeValue
	if order.Meta != nil {
		meta = make(map[string]types.AttributeValue)
		for k, v := range order.Meta {
			meta[k] = &types.AttributeValueMemberS{Value: v}
		}
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: tenantKey},
			"SK":         &types.AttributeValueMemberS{Value: orderKey},
			"Type":       &types.AttributeValueMemberS{Value: "Order"},
			"CustomerID": &types.AttributeValueMemberS{Value: order.CustomerID},
			"Meta":       &types.AttributeValueMemberM{Value: meta},
			"OrderItems": &types.AttributeValueMemberL{Value: orderItems},
		},
	}

	_, err := repo.client.PutItem(context.TODO(), input)
	if err != nil {
		return errors.Wrap(err, "failed to put item")
	}

	return nil
}
