package square

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/square"
	"os"
)

type PaymentsRepositoryConfig struct {
	TableName string
}

func NewPaymentsRepositoryConfig() PaymentsRepositoryConfig {
	return PaymentsRepositoryConfig{
		TableName: os.Getenv(SquarePaymentsTableName),
	}
}

type PaymentsRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewPaymentsRepository(config PaymentsRepositoryConfig, client *dynamodb.Client) square.PaymentsRepository {
	return PaymentsRepository{
		tableName: config.TableName,
		client:    client,
	}
}

func (repo PaymentsRepository) Create(payment square.Payment) error {
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
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (repo PaymentsRepository) FindByID(squarePaymentID string) (square.Payment, error) {
	var payment = square.Payment{}

	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			SquarePaymentID: &types.AttributeValueMemberS{Value: squarePaymentID},
		},
		TableName: aws.String(repo.tableName),
	}

	result, err := repo.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return payment, fmt.Errorf("failed to get payment with id '%s': %w", squarePaymentID, err)
	}

	err = attributevalue.UnmarshalMap(result.Item, &payment)
	if err != nil {
		return payment, fmt.Errorf("failed to unmarshal payment: %w", err)
	}

	return payment, nil
}