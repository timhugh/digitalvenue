package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/square/db"
	"os"
)

type PaymentsRepositoryConfig struct {
	TableName string
}

func NewPaymentsRepositoryConfig() PaymentsRepositoryConfig {
	return PaymentsRepositoryConfig{
		TableName: os.Getenv(PaymentsTableName),
	}
}

type PaymentsRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewPaymentsRepository(config PaymentsRepositoryConfig, client *dynamodb.Client) db.PaymentsRepository {
	return PaymentsRepository{
		tableName: config.TableName,
		client:    client,
	}
}

func (repo PaymentsRepository) CreatePayment(payment db.Payment) error {
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
