package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
	"os"
)

const (
	SquarePaymentID  = "SquarePaymentID"
	SquareMerchantID = "SquareMerchantID"
	SquareOrderID    = "SquareOrderID"
)

type PaymentsRepositoryConfig struct {
	TableName string
}

func NewPaymentsRepositoryConfig() PaymentsRepositoryConfig {
	return PaymentsRepositoryConfig{
		TableName: os.Getenv("PAYMENTS_TABLE"),
	}
}

type PaymentsRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewPaymentsRepository(config PaymentsRepositoryConfig, client *dynamodb.Client) PaymentsRepository {
	return PaymentsRepository{
		tableName: config.TableName,
		client:    client,
	}
}

func (p PaymentsRepository) CreatePayment(payment core.Payment) error {
	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			SquarePaymentID:  &types.AttributeValueMemberS{Value: payment.SquarePaymentID},
			SquareMerchantID: &types.AttributeValueMemberS{Value: payment.SquareMerchantID},
			SquareOrderID:    &types.AttributeValueMemberS{Value: payment.SquareOrderID},
		},
		TableName: aws.String(p.tableName),
	}

	_, err := p.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create payment")
		return fmt.Errorf("failed to create payment")
	}

	return nil
}
