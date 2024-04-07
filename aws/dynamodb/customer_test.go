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

const customersTableName = "customers-test"

func initCustomerRepositoryTest(t *testing.T) *is.I {
	is := is.New(t)
	mock.SetUp(t)

	err := os.Setenv(CustomersTableNameKey, customersTableName)
	is.NoErr(err)

	return is
}

func TestCustomerRepository_NewCustomerRepository_RequiresTableName(t *testing.T) {
	is := initCustomerRepositoryTest(t)

	err := os.Unsetenv(CustomersTableNameKey)
	is.NoErr(err)

	_, err = NewCustomerRepository(mock.Mock[Client]())
	is.Equal(err.Error(), "missing required environment variable CUSTOMERS_TABLE_NAME")
}

func TestCustomerRepository_PutCustomer(t *testing.T) {
	is := initCustomerRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()

	mockIDGenerator := mock.Mock[core.IDGenerator]()
	mock.WhenSingle(mockIDGenerator.GenerateID()).ThenReturn(test.CustomerID)

	repo := CustomerRepository{
		client:      mockDynamoDBClient,
		idGenerator: mockIDGenerator,
		tableName:   customersTableName,
	}

	customer := test.NewCustomer()
	customer.Meta = map[string]string{
		"ExampleKey": "exampleValue",
	}

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).ThenReturn(nil, nil)

	customerID, err := repo.PutCustomer(customer)
	is.NoErr(err)
	is.Equal(customerID, test.CustomerID)

	putItemInput := putItemInputCaptor.Last()
	expectedPutItemInput := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			CustomerID: &types.AttributeValueMemberS{Value: test.CustomerID},
			FirstName:  &types.AttributeValueMemberS{Value: test.CustomerFirstName},
			LastName:   &types.AttributeValueMemberS{Value: test.CustomerLastName},
			Email:      &types.AttributeValueMemberS{Value: test.CustomerEmail},
			Phone:      &types.AttributeValueMemberS{Value: test.CustomerPhone},
			Meta: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"ExampleKey": &types.AttributeValueMemberS{Value: "exampleValue"},
			}},
		},
		TableName: aws.String(customersTableName),
	}
	err = test.Diff(expectedPutItemInput, putItemInput)
	is.NoErr(err)
}

func TestCustomerRepository_PutCustomer_ExistingCustomer(t *testing.T) {
	is := initCustomerRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()

	repo, err := NewCustomerRepository(mockDynamoDBClient)
	is.NoErr(err)

	customer := test.NewCustomer()
	customer.CustomerID = "existing_customer"

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).ThenReturn(nil, nil)

	customerID, err := repo.PutCustomer(customer)
	is.NoErr(err)
	is.Equal(customerID, "existing_customer")

	putItemInput := putItemInputCaptor.Last()
	is.Equal("existing_customer", putItemInput.Item["CustomerID"].(*types.AttributeValueMemberS).Value)
}

func TestCustomerRepository_PutCustomer_PutError(t *testing.T) {
	is := initCustomerRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()

	repo, err := NewCustomerRepository(mockDynamoDBClient)
	is.NoErr(err)

	thrownError := fmt.Errorf("some error from dynamodb")

	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), mock.Any[*dynamodb.PutItemInput]())).ThenReturn(nil, thrownError)

	_, err = repo.PutCustomer(core.Customer{})
	is.True(errors.Is(err, thrownError))
}
