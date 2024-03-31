package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/db"
)

const (
	SquarePaymentID  = "SquarePaymentID"
	SquareMerchantID = "SquareMerchantID"
	SquareOrderID    = "SquareOrderID"
)

type PaymentsRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewPaymentsRepository(tableName string) (db.PaymentsRepository, error) {
	client, err := Connect()
	if err != nil {
		return nil, err
	}

	return &PaymentsRepository{
		tableName: tableName,
		client:    client,
	}, nil
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
