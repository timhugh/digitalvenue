package dv_dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/square/squaretest"
	"github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func squarePaymentGetItemOutput() *dynamodb.GetItemOutput {
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"PK":            &types.AttributeValueMemberS{Value: "SquareMerchant#" + squaretest.SquareMerchantID},
			"SK":            &types.AttributeValueMemberS{Value: "SquarePayment#" + squaretest.SquarePaymentID},
			"Type":          &types.AttributeValueMemberS{Value: "SquarePayment"},
			"TenantID":      &types.AttributeValueMemberS{Value: "Tenant#" + test.TenantID},
			"SquareOrderID": &types.AttributeValueMemberS{Value: squaretest.SquareOrderID},
		},
	}
}

func squarePaymentGetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "SquareMerchant#" + squaretest.SquareMerchantID},
			"SK": &types.AttributeValueMemberS{Value: "SquarePayment#" + squaretest.SquarePaymentID},
		},
	}
}

func squarePaymentPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK":            &types.AttributeValueMemberS{Value: "SquareMerchant#" + squaretest.SquareMerchantID},
			"SK":            &types.AttributeValueMemberS{Value: "SquarePayment#" + squaretest.SquarePaymentID},
			"Type":          &types.AttributeValueMemberS{Value: "SquarePayment"},
			"SquareOrderID": &types.AttributeValueMemberS{Value: squaretest.SquareOrderID},
		},
	}
}

func TestRepository_GetSquarePayment_BasicSuccess(t *testing.T) {
	repo, client := initRepositoryTest(t)

	getItemInputCaptor := mock.Captor[*dynamodb.GetItemInput]()
	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), getItemInputCaptor.Capture())).
		ThenReturn(squarePaymentGetItemOutput(), nil)

	payment, err := repo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := test.Diff(squaretest.NewSquarePayment(), payment); err != nil {
		t.Error(err)
	}

	if err := test.Diff(squarePaymentGetItemInput(), getItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_GetSquarePayment_GetItemError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	inducedError := errors.New("induced error")

	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(nil, inducedError)

	_, err := repo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	if !errors.Is(err, inducedError) {
		t.Errorf("expected error %v, got %v", inducedError, err)
	}
}

func TestRepository_GetSquarePayment_NoItemError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(&dynamodb.GetItemOutput{}, nil)

	payment, err := repo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if payment != nil {
		t.Errorf("expected nil, got %+v", payment)
	}
}

func TestRepository_GetSquarePayment_IncorrectItemTypeError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	wrongItem := squarePaymentGetItemOutput()
	wrongItem.Item["Type"] = &types.AttributeValueMemberS{Value: "SquareMerchant"}

	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(wrongItem, nil)

	_, err := repo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestRepository_PutSquarePayment_BasicSuccess(t *testing.T) {
	repo, client := initRepositoryTest(t)

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(client.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).
		ThenReturn(nil, nil)

	payment := squaretest.NewSquarePayment()
	err := repo.PutSquarePayment(payment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := test.Diff(squarePaymentPutItemInput(), putItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_PutSquarePayment_PutItemError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	inducedError := errors.New("induced error")

	mock.WhenDouble(client.PutItem(mock.Any[context.Context](), mock.Any[*dynamodb.PutItemInput]())).
		ThenReturn(nil, inducedError)

	err := repo.PutSquarePayment(squaretest.NewSquarePayment())
	if !errors.Is(err, inducedError) {
		t.Errorf("expected error %v, got %v", inducedError, err)
	}
}
