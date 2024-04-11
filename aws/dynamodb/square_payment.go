package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/square"
)

type squarePayment struct {
	PK            string
	SK            string
	Type          string
	TenantID      string
	SquareOrderID string
}

func (repo *Repository) GetSquarePayment(squareMerchantID string, squarePaymentID string) (*square.Payment, error) {
	pk := PrefixID("SquareMerchant", squareMerchantID)
	sk := PrefixID("SquarePayment", squarePaymentID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	item := squarePayment{}
	err := repo.getItem("SquarePayment", input, &item)
	if err != nil {
		return nil, err
	}

	return &square.Payment{
		SquarePaymentID:  squarePaymentID,
		SquareMerchantID: squareMerchantID,
		SquareOrderID:    item.SquareOrderID,
	}, nil
}

func (repo *Repository) PutSquarePayment(payment *square.Payment) error {
	pk := PrefixID("SquareMerchant", payment.SquareMerchantID)
	sk := PrefixID("SquarePayment", payment.SquarePaymentID)
	input := &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item: map[string]types.AttributeValue{
			"PK":            &types.AttributeValueMemberS{Value: pk},
			"SK":            &types.AttributeValueMemberS{Value: sk},
			"Type":          &types.AttributeValueMemberS{Value: "SquarePayment"},
			"SquareOrderID": &types.AttributeValueMemberS{Value: payment.SquareOrderID},
		},
	}

	_, err := repo.client.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
