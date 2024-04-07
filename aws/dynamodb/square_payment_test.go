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

const squarePaymentTableName = "square-payments-test"

func initSquarePaymentRepositoryTest(t *testing.T) *is.I {
	is := is.New(t)
	mock.SetUp(t)

	err := os.Setenv(SquarePaymentsTableNameKey, squarePaymentTableName)
	is.NoErr(err)

	return is
}

func TestSquarePaymentRepository_NewSquarePaymentRepository_RequiresTableName(t *testing.T) {
	is := initSquarePaymentRepositoryTest(t)

	err := os.Unsetenv(SquarePaymentsTableNameKey)
	is.NoErr(err)

	_, err = NewSquarePaymentRepository(mock.Mock[Client]())
	is.Equal(err.Error(), "missing required environment variable SQUARE_PAYMENTS_TABLE_NAME")
}

func TestSquarePaymentRepository_GetSquarePayment(t *testing.T) {
	is := initSquarePaymentRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()
	getItemInputCaptor := mock.Captor[*dynamodb.GetItemInput]()
	getItemOutput := dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			SquarePaymentID:  &types.AttributeValueMemberS{Value: squaretest.SquarePaymentID},
			SquareMerchantID: &types.AttributeValueMemberS{Value: squaretest.SquareMerchantID},
			SquareOrderID:    &types.AttributeValueMemberS{Value: squaretest.SquareOrderID},
		},
	}
	mock.WhenDouble(mockDynamoDBClient.GetItem(mock.Any[context.Context](), getItemInputCaptor.Capture())).ThenReturn(&getItemOutput, nil)

	repo, err := NewSquarePaymentRepository(mockDynamoDBClient)
	is.NoErr(err)

	payment, err := repo.GetSquarePayment(squaretest.SquarePaymentID)
	is.NoErr(err)

	getItemInput := getItemInputCaptor.Last()
	expectedGetItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			SquarePaymentID: &types.AttributeValueMemberS{Value: squaretest.SquarePaymentID},
		},
		TableName: aws.String(squarePaymentTableName),
	}
	err = test.Diff(expectedGetItemInput, getItemInput)
	is.NoErr(err)

	expectedPayment := squaretest.NewSquarePayment()
	err = test.Diff(expectedPayment, payment)
	is.NoErr(err)
}

func TestSquarePaymentRepository_GetSquarePayment_GetItemError(t *testing.T) {
	is := initSquarePaymentRepositoryTest(t)

	thrownError := errors.New("some client error")

	mockDynamoDBClient := mock.Mock[Client]()
	mock.WhenDouble(mockDynamoDBClient.GetItem(mock.Any[context.Context](), mock.Any[*dynamodb.GetItemInput]())).ThenReturn(nil, thrownError)

	repo, err := NewSquarePaymentRepository(mockDynamoDBClient)
	is.NoErr(err)

	_, err = repo.GetSquarePayment(squaretest.SquarePaymentID)
	is.True(errors.Is(err, thrownError))
}

func TestSquarePaymentRepository_PutSquarePayment(t *testing.T) {
	is := initSquarePaymentRepositoryTest(t)

	mockDynamoDBClient := mock.Mock[Client]()
	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).ThenReturn(nil, nil)

	repo, err := NewSquarePaymentRepository(mockDynamoDBClient)
	is.NoErr(err)

	payment := squaretest.NewSquarePayment()

	err = repo.PutSquarePayment(payment)
	is.NoErr(err)

	putItemInput := putItemInputCaptor.Last()
	expectedPutItemInput := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			SquarePaymentID:  &types.AttributeValueMemberS{Value: squaretest.SquarePaymentID},
			SquareMerchantID: &types.AttributeValueMemberS{Value: squaretest.SquareMerchantID},
			SquareOrderID:    &types.AttributeValueMemberS{Value: squaretest.SquareOrderID},
		},
		TableName: aws.String(squarePaymentTableName),
	}
	err = test.Diff(expectedPutItemInput, putItemInput)
	is.NoErr(err)
}

func TestSquarePaymentRepository_PutSquarePayment_PutItemError(t *testing.T) {
	is := initSquarePaymentRepositoryTest(t)

	thrownError := errors.New("some client error")

	mockDynamoDBClient := mock.Mock[Client]()
	mock.WhenDouble(mockDynamoDBClient.PutItem(mock.Any[context.Context](), mock.Any[*dynamodb.PutItemInput]())).ThenReturn(nil, thrownError)

	repo, err := NewSquarePaymentRepository(mockDynamoDBClient)
	is.NoErr(err)

	err = repo.PutSquarePayment(square.Payment{})
	is.Equal(thrownError, err)
}
