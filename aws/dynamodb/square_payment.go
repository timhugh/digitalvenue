package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square"
)

type SquarePaymentRepository struct {
	client      Client
	idGenerator core.IDGenerator
	tableName   string
}

func NewSquarePaymentRepository(client Client) (*SquarePaymentRepository, error) {
	tableName, err := core.RequireEnv(SquarePaymentsTableNameKey)
	if err != nil {
		return nil, err
	}

	return &SquarePaymentRepository{
		client:      client,
		idGenerator: core.NewIDGenerator(),
		tableName:   tableName,
	}, nil
}

func (repo *SquarePaymentRepository) PutSquarePayment(payment square.Payment) error {
	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			SquarePaymentID:  &types.AttributeValueMemberS{Value: payment.SquarePaymentID},
			SquareMerchantID: &types.AttributeValueMemberS{Value: payment.SquareMerchantID},
			SquareOrderID:    &types.AttributeValueMemberS{Value: payment.SquareOrderID},
		},
		TableName: aws.String(repo.tableName),
	}

	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SquarePaymentRepository) GetSquarePayment(squarePaymentID string) (square.Payment, error) {
	var payment = square.Payment{}

	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			SquarePaymentID: &types.AttributeValueMemberS{Value: squarePaymentID},
		},
		TableName: aws.String(repo.tableName),
	}

	result, err := repo.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return payment, err
	}

	_ = attributevalue.UnmarshalMap(result.Item, &payment)

	return payment, nil
}
