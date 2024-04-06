package dynamodb

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"github.com/timhugh/digitalvenue/test"
	"os"
	"testing"
)

const squareMerchantTableName = "square-merchants-test"

func initSquareMerchantRepositoryTest(t *testing.T) *is.I {
	is := is.New(t)
	mock.SetUp(t)

	err := os.Setenv(SquareMerchantsTableNameKey, squareMerchantTableName)
	is.NoErr(err)

	return is
}

func TestSquareMerchantRepository_NewSquareMerchantRepository_RequiresTableName(t *testing.T) {
	is := initSquareMerchantRepositoryTest(t)

	err := os.Unsetenv(SquareMerchantsTableNameKey)
	is.NoErr(err)

	_, err = NewSquareMerchantRepository(mock.Mock[Client]())
	is.Equal(err.Error(), "missing required environment variable SQUARE_MERCHANTS_TABLE_NAME")
}

func TestSquareMerchantRepository_GetSquareMerchant(t *testing.T) {
	is := initSquareMerchantRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()
	getItemInputCaptor := mock.Captor[*dynamodb.GetItemInput]()
	getItemOutput := dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			SquareMerchantID:          &types.AttributeValueMemberS{Value: squaretest.SquareMerchantID},
			SquareWebhookSignatureKey: &types.AttributeValueMemberS{Value: squaretest.SquareWebhookSignatureKey},
			SquareAPIToken:            &types.AttributeValueMemberS{Value: squaretest.SquareAPIToken},
		},
	}
	mock.WhenDouble(mockDynamoDBClient.GetItem(mock.Any[context.Context](), getItemInputCaptor.Capture())).ThenReturn(&getItemOutput, nil)

	repo, err := NewSquareMerchantRepository(mockDynamoDBClient)
	is.NoErr(err)

	merchant, err := repo.GetSquareMerchant(squaretest.SquareMerchantID)
	is.NoErr(err)

	getItemInput := getItemInputCaptor.Last()
	expectedGetItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			SquareMerchantID: &types.AttributeValueMemberS{Value: squaretest.SquareMerchantID},
		},
		TableName: aws.String(squareMerchantTableName),
	}
	err = test.Diff(expectedGetItemInput, getItemInput)
	is.NoErr(err)

	expectedMerchant := squaretest.NewSquareMerchant()
	err = test.Diff(expectedMerchant, merchant)
	is.NoErr(err)
}

func TestSquareMerchantRepository_GetSquareMerchant_GetItemError(t *testing.T) {
	is := initSquareMerchantRepositoryTest(t)

	thrownError := errors.New("some client error")

	mockDynamoDBClient := mock.Mock[Client]()
	mock.WhenDouble(mockDynamoDBClient.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).ThenReturn(nil, thrownError)

	repo, err := NewSquareMerchantRepository(mockDynamoDBClient)
	is.NoErr(err)

	_, err = repo.GetSquareMerchant(squaretest.SquareMerchantID)
	is.Equal(err, thrownError)
}

func TestSquareMerchantRepository_PutSquareMerchant(t *testing.T) {
	is := initSquareMerchantRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()
	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).ThenReturn(nil, nil)

	repo, err := NewSquareMerchantRepository(mockDynamoDBClient)
	is.NoErr(err)

	merchant := squaretest.NewSquareMerchant()

	err = repo.PutSquareMerchant(merchant)
	is.NoErr(err)

	putItemInput := putItemInputCaptor.Last()
	expectedPutItemInput := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			SquareMerchantID:          &types.AttributeValueMemberS{Value: merchant.SquareMerchantID},
			SquareWebhookSignatureKey: &types.AttributeValueMemberS{Value: merchant.SquareWebhookSignatureKey},
			SquareAPIToken:            &types.AttributeValueMemberS{Value: merchant.SquareAPIToken},
		},
		TableName: aws.String(squareMerchantTableName),
	}
	err = test.Diff(expectedPutItemInput, putItemInput)
	is.NoErr(err)
}

func TestSquareMerchantRepository_PutSquareMerchant_PutItemError(t *testing.T) {
	is := initSquareMerchantRepositoryTest(t)

	thrownError := errors.New("some client error")

	mockDynamoDBClient := mock.Mock[Client]()
	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), mock.Any[*dynamodb.PutItemInput]())).ThenReturn(nil, thrownError)

	repo, err := NewSquareMerchantRepository(mockDynamoDBClient)
	is.NoErr(err)

	err = repo.PutSquareMerchant(square.Merchant{})
	is.Equal(err, thrownError)
}
