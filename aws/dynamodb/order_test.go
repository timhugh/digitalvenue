package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/test"
	"os"
	"testing"
)

const ordersTableName = "orders-test"

func initOrderRepositoryTest(t *testing.T) *is.I {
	is := is.New(t)
	mock.SetUp(t)

	err := os.Setenv(OrdersTableNameKey, ordersTableName)
	is.NoErr(err)

	return is
}

func TestOrderRepository_NewOrderRepository_RequiresTableName(t *testing.T) {
	is := initOrderRepositoryTest(t)

	err := os.Unsetenv(OrdersTableNameKey)
	is.NoErr(err)

	_, err = NewOrderRepository(mock.Mock[Client]())
	is.Equal(err.Error(), "missing required environment variable ORDERS_TABLE_NAME")
}

func TestOrderRepository_PutOrder(t *testing.T) {
	is := initOrderRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()

	mockIDGenerator := mock.Mock[core.IDGenerator]()
	mock.WhenSingle(mockIDGenerator.GenerateID()).ThenReturn(test.OrderID)

	repo := OrderRepository{
		client:      mockDynamoDBClient,
		idGenerator: mockIDGenerator,
		tableName:   ordersTableName,
	}

	order := test.NewOrder()
	order.CustomerID = test.CustomerID
	order.Meta = map[string]string{
		"ExampleKey": "exampleValue",
	}
	order.Items[0].Meta = map[string]string{
		"ExampleKey2": "exampleValue2",
	}

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).ThenReturn(nil, nil)

	orderID, err := repo.PutOrder(order)
	is.NoErr(err)
	is.Equal(orderID, test.OrderID)

	putItemInput := putItemInputCaptor.Last()
	expectedPutItemInput := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			OrderID:    &types.AttributeValueMemberS{Value: test.OrderID},
			TenantID:   &types.AttributeValueMemberS{Value: ""},
			CustomerID: &types.AttributeValueMemberS{Value: test.CustomerID},
			Items: &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						ItemID: &types.AttributeValueMemberS{Value: test.OrderID}, // TODO: not sure how to get the id generator mock to return sequential values
						Name:   &types.AttributeValueMemberS{Value: test.ItemName1},
						Meta: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
							"ExampleKey2": &types.AttributeValueMemberS{Value: "exampleValue2"},
						}},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						ItemID: &types.AttributeValueMemberS{Value: test.OrderID},
						Name:   &types.AttributeValueMemberS{Value: test.ItemName2},
						Meta:   &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						ItemID: &types.AttributeValueMemberS{Value: test.OrderID},
						Name:   &types.AttributeValueMemberS{Value: test.ItemName2},
						Meta:   &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
					},
				},
			}},
			Meta: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"ExampleKey": &types.AttributeValueMemberS{Value: "exampleValue"},
			}},
		},
		TableName: aws.String(ordersTableName),
	}
	err = test.Diff(expectedPutItemInput, putItemInput)
	is.NoErr(err)
}

func TestOrderRepository_PutOrder_ExistingOrder(t *testing.T) {
	is := initOrderRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()

	repo, err := NewOrderRepository(mockDynamoDBClient)
	is.NoErr(err)

	order := test.NewOrder()
	order.OrderID = "existing_order"
	order.Items[0].ItemID = "existing_item1"
	order.Items[1].ItemID = "existing_item2"
	order.Items[2].ItemID = "existing_item3"

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).ThenReturn(nil, nil)

	orderID, err := repo.PutOrder(order)
	is.NoErr(err)
	is.Equal(orderID, "existing_order")

	putItemInput := putItemInputCaptor.Last()
	is.Equal("existing_order", putItemInput.Item[OrderID].(*types.AttributeValueMemberS).Value)
}

func TestOrderRepository_PutOrder_PutError(t *testing.T) {
	is := initOrderRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()

	repo, err := NewOrderRepository(mockDynamoDBClient)
	is.NoErr(err)

	thrownError := fmt.Errorf("some error from dynamodb")

	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), mock.Any[*dynamodb.PutItemInput]())).ThenReturn(nil, thrownError)

	_, err = repo.PutOrder(core.Order{})
	is.True(errors.Is(err, thrownError))
}
