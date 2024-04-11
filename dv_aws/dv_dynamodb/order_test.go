package dv_dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/test"
	"testing"
)

func orderPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: "Tenant#" + test.TenantID},
			"SK":         &types.AttributeValueMemberS{Value: "Order#" + test.OrderID},
			"Type":       &types.AttributeValueMemberS{Value: "Order"},
			"CustomerID": &types.AttributeValueMemberS{Value: test.CustomerID},
			"Meta":       &types.AttributeValueMemberM{Value: nil},
			"OrderItems": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ItemID": &types.AttributeValueMemberS{Value: test.ItemID1},
						"Name":   &types.AttributeValueMemberS{Value: test.ItemName1},
						"Meta":   &types.AttributeValueMemberM{Value: nil},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ItemID": &types.AttributeValueMemberS{Value: test.ItemID2},
						"Name":   &types.AttributeValueMemberS{Value: test.ItemName2},
						"Meta":   &types.AttributeValueMemberM{Value: nil},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ItemID": &types.AttributeValueMemberS{Value: test.ItemID3},
						"Name":   &types.AttributeValueMemberS{Value: test.ItemName2},
						"Meta":   &types.AttributeValueMemberM{Value: nil},
					},
				},
			}},
		},
	}
}

func TestRepository_PutOrder_Success(t *testing.T) {
	repo, client := initRepositoryTest(t)

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(client.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).
		ThenReturn(nil, nil)

	order := test.NewOrder()
	err := repo.PutOrder(order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := test.Diff(orderPutItemInput(), putItemInputCaptor.Last()); err != nil {
		t.Fatalf("unexpected input: %v", err)
	}
}
