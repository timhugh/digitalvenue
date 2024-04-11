package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"github.com/timhugh/digitalvenue/test"
	"testing"
)

func squareMerchantGetItemOutput() *dynamodb.GetItemOutput {
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"PK":                        &types.AttributeValueMemberS{Value: "SquareMerchant#" + squaretest.SquareMerchantID},
			"SK":                        &types.AttributeValueMemberS{Value: "SquareMerchant#" + squaretest.SquareMerchantID},
			"Type":                      &types.AttributeValueMemberS{Value: "SquareMerchant"},
			"TenantID":                  &types.AttributeValueMemberS{Value: "Tenant#" + test.TenantID},
			"Name":                      &types.AttributeValueMemberS{Value: test.TenantName},
			"SquareAPIToken":            &types.AttributeValueMemberS{Value: squaretest.SquareAPIToken},
			"SquareWebhookSignatureKey": &types.AttributeValueMemberS{Value: squaretest.SquareWebhookSignatureKey},
		},
	}
}

func squareMerchantGetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "SquareMerchant#" + squaretest.SquareMerchantID},
			"SK": &types.AttributeValueMemberS{Value: "SquareMerchant#" + squaretest.SquareMerchantID},
		},
	}
}

func TestRepository_GetSquareMerchant_BasicSuccess(t *testing.T) {
	repo, client := initRepositoryTest(t)

	getItemInputCaptor := mock.Captor[*dynamodb.GetItemInput]()
	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), getItemInputCaptor.Capture())).
		ThenReturn(squareMerchantGetItemOutput(), nil)

	merchant, err := repo.GetSquareMerchant(squaretest.SquareMerchantID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := test.Diff(squaretest.NewSquareMerchant(), merchant); err != nil {
		t.Error(err)
	}

	if err := test.Diff(squareMerchantGetItemInput(), getItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_GetSquareMerchant_GetItemError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	inducedError := errors.New("induced error")

	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(nil, inducedError)

	_, err := repo.GetSquareMerchant(squaretest.SquareMerchantID)
	if !errors.Is(err, inducedError) {
		t.Errorf("expected error %v, got %v", inducedError, err)
	}
}

func TestRepository_GetSquareMerchant_NoItemError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(&dynamodb.GetItemOutput{}, nil)

	merchant, err := repo.GetSquareMerchant(squaretest.SquareMerchantID)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if merchant != nil {
		t.Errorf("expected nil, got %+v", merchant)
	}
}

func TestRepository_GetSquareMerchant_IncorrectItemTypeError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	wrongItem := &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"Type": &types.AttributeValueMemberS{Value: "NotSquareMerchant"},
		},
	}

	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(wrongItem, nil)

	_, err := repo.GetSquareMerchant(squaretest.SquareMerchantID)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestRepository_GetSquareMerchant_InvalidTenantIDError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	outputWithInvalidTenantID := squareMerchantGetItemOutput()
	outputWithInvalidTenantID.Item["TenantID"] = &types.AttributeValueMemberS{Value: "InvalidTenantID"}

	mock.WhenDouble(client.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(outputWithInvalidTenantID, nil)

	_, err := repo.GetSquareMerchant(squaretest.SquareMerchantID)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
