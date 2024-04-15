package dv_dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ovechkin-dm/mockio/mock"
	test2 "github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func orderPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: "Tenant#" + test2.TenantID},
			"SK":         &types.AttributeValueMemberS{Value: "Order#" + test2.OrderID},
			"Type":       &types.AttributeValueMemberS{Value: "Order"},
			"CustomerID": &types.AttributeValueMemberS{Value: test2.CustomerID},
			"Meta":       &types.AttributeValueMemberM{Value: nil},
			"OrderItems": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ItemID": &types.AttributeValueMemberS{Value: test2.ItemID1},
						"Name":   &types.AttributeValueMemberS{Value: test2.ItemName1},
						"Meta":   &types.AttributeValueMemberM{Value: nil},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ItemID": &types.AttributeValueMemberS{Value: test2.ItemID2},
						"Name":   &types.AttributeValueMemberS{Value: test2.ItemName2},
						"Meta":   &types.AttributeValueMemberM{Value: nil},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ItemID": &types.AttributeValueMemberS{Value: test2.ItemID3},
						"Name":   &types.AttributeValueMemberS{Value: test2.ItemName2},
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

	order := test2.NewOrder()
	err := repo.PutOrder(order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := test2.Diff(orderPutItemInput(), putItemInputCaptor.Last()); err != nil {
		t.Fatalf("unexpected input: %v", err)
	}
}
