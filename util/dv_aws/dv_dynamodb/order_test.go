package dv_dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func orderAttributes() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK":         &types.AttributeValueMemberS{Value: "Tenant#" + test.TenantID},
		"SK":         &types.AttributeValueMemberS{Value: "Order#" + test.OrderID},
		"Type":       &types.AttributeValueMemberS{Value: "Order"},
		"CustomerID": &types.AttributeValueMemberS{Value: test.CustomerID},
		"Meta":       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
		"OrderItems": &types.AttributeValueMemberL{Value: []types.AttributeValue{
			&types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"ItemID": &types.AttributeValueMemberS{Value: test.ItemID1},
					"Name":   &types.AttributeValueMemberS{Value: test.ItemName1},
					"Meta":   &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				},
			},
			&types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"ItemID": &types.AttributeValueMemberS{Value: test.ItemID2},
					"Name":   &types.AttributeValueMemberS{Value: test.ItemName2},
					"Meta":   &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				},
			},
			&types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"ItemID": &types.AttributeValueMemberS{Value: test.ItemID3},
					"Name":   &types.AttributeValueMemberS{Value: test.ItemName2},
					"Meta":   &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				},
			},
		}},
	}
}

func TestRepository_PutOrder_Success(t *testing.T) {
	repo, client := initRepositoryTest(t)

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(client.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).ThenReturn(nil, nil)

	order := test.NewOrder()
	err := repo.PutOrder(order)
	if err != nil {
		t.Fatal(err)
	}

	orderPutItemInput := dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      orderAttributes(),
	}
	if err := test.Diff(&orderPutItemInput, putItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_GetOrder_Success(t *testing.T) {
	repo, client := initRepositoryTest(t)

	getItemOutput := dynamodb.GetItemOutput{
		Item: orderAttributes(),
	}

	getItemInputCaptor := mock.Captor[*dynamodb.GetItemInput]()
	mock.When(client.GetItem(mock.Any[context.Context](), getItemInputCaptor.Capture())).ThenReturn(&getItemOutput, nil)

	order, err := repo.GetOrder(test.TenantID, test.OrderID)
	if err != nil {
		t.Fatal(err)
	}

	expectedOrder := test.NewOrder()
	if err := test.Diff(expectedOrder, order); err != nil {
		t.Error(err)
	}

	expectedGetItemInput := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "Tenant#" + test.TenantID},
			"SK": &types.AttributeValueMemberS{Value: "Order#" + test.OrderID},
		},
	}
	if err := test.Diff(&expectedGetItemInput, getItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}
