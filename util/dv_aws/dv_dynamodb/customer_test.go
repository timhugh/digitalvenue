package dv_dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/pkg/errors"
	test2 "github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func customerGetItemOutput() *dynamodb.GetItemOutput {
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: "Tenant#" + test2.TenantID},
			"SK":         &types.AttributeValueMemberS{Value: "Customer#" + test2.CustomerID},
			"Type":       &types.AttributeValueMemberS{Value: "Customer"},
			"CustomerID": &types.AttributeValueMemberS{Value: test2.CustomerID},
			"Name":       &types.AttributeValueMemberS{Value: test2.CustomerName},
			"Email":      &types.AttributeValueMemberS{Value: test2.CustomerEmail},
			"Phone":      &types.AttributeValueMemberS{Value: test2.CustomerPhone},
		},
	}
}

func customerGetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "Tenant#" + test2.TenantID},
			"SK": &types.AttributeValueMemberS{Value: "Customer#" + test2.CustomerID},
		},
	}
}

func customerPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: "Tenant#" + test2.TenantID},
			"SK":         &types.AttributeValueMemberS{Value: "Customer#" + test2.CustomerID},
			"Type":       &types.AttributeValueMemberS{Value: "Customer"},
			"CustomerID": &types.AttributeValueMemberS{Value: test2.CustomerID},
			"Name":       &types.AttributeValueMemberS{Value: test2.CustomerName},
			"Email":      &types.AttributeValueMemberS{Value: test2.CustomerEmail},
			"Phone":      &types.AttributeValueMemberS{Value: test2.CustomerPhone},
			"Meta":       &types.AttributeValueMemberM{Value: nil},
		},
	}
}

func TestRepository_GetCustomer_Success(t *testing.T) {
	repo, client := initRepositoryTest(t)

	getItemInputCaptor := mock.Captor[*dynamodb.GetItemInput]()
	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), getItemInputCaptor.Capture())).
		ThenReturn(customerGetItemOutput(), nil)

	actualCustomer, err := repo.GetCustomer(test2.TenantID, test2.CustomerID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := test2.Diff(test2.NewCustomer(), actualCustomer); err != nil {
		t.Error(err)
	}

	if err := test2.Diff(customerGetItemInput(), getItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_GetCustomer_SuccessWithMetadata(t *testing.T) {
	repo, client := initRepositoryTest(t)

	getItemOutput := customerGetItemOutput()
	getItemOutput.Item["Meta"] = &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"ExampleKey": &types.AttributeValueMemberS{Value: "ExampleValue"},
		},
	}

	getItemInputCaptor := mock.Captor[*dynamodb.GetItemInput]()
	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), getItemInputCaptor.Capture())).
		ThenReturn(getItemOutput, nil)

	actualCustomer, err := repo.GetCustomer(test2.TenantID, test2.CustomerID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCustomer := test2.NewCustomer()
	expectedCustomer.Meta = map[string]string{
		"ExampleKey": "ExampleValue",
	}

	if err := test2.Diff(expectedCustomer, actualCustomer); err != nil {
		t.Error(err)
	}

	if err := test2.Diff(customerGetItemInput(), getItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_GetCustomer_GetItemError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	inducedError := errors.New("induced error")

	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(nil, inducedError)

	_, err := repo.GetCustomer(test2.TenantID, test2.CustomerID)
	if !errors.Is(err, inducedError) {
		t.Errorf("expected error %v, got %v", inducedError, err)
	}
}

func TestRepository_PutCustomer_Success(t *testing.T) {
	repo, client := initRepositoryTest(t)

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(client.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).
		ThenReturn(nil, nil)

	err := repo.PutCustomer(test2.NewCustomer())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := test2.Diff(customerPutItemInput(), putItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_PutCustomer_SuccessWithMetadata(t *testing.T) {
	repo, client := initRepositoryTest(t)

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(client.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).
		ThenReturn(nil, nil)

	customer := test2.NewCustomer()
	customer.Meta = map[string]string{
		"ExampleKey": "ExampleValue",
	}

	err := repo.PutCustomer(customer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedInput := customerPutItemInput()
	expectedInput.Item["Meta"] = &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"ExampleKey": &types.AttributeValueMemberS{Value: "ExampleValue"},
		},
	}
	if err := test2.Diff(expectedInput, putItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_PutCustomer_PutItemError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	inducedError := errors.New("induced error")

	mock.WhenDouble(client.PutItem(mock.Any[context.Context](), mock.Any[*dynamodb.PutItemInput]())).
		ThenReturn(nil, inducedError)

	err := repo.PutCustomer(test2.NewCustomer())
	if !errors.Is(err, inducedError) {
		t.Errorf("expected error %v, got %v", inducedError, err)
	}
}
