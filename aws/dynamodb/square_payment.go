package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
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
	pk := "SquareMerchant#" + squareMerchantID
	sk := "SquarePayment#" + squarePaymentID
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

	tenantID, err := removeIDPrefix(item.TenantID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid tenant ID")
	}

	return &square.Payment{
		SquarePaymentID:  squarePaymentID,
		SquareMerchantID: squareMerchantID,
		SquareOrderID:    item.SquareOrderID,
		TenantID:         tenantID,
	}, nil
}

func (repo *Repository) PutSquarePayment(payment *square.Payment) error {
	pk := "SquareMerchant#" + payment.SquareMerchantID
	sk := "SquarePayment#" + payment.SquarePaymentID
	tenantID := "Tenant#" + payment.TenantID
	input := &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item: map[string]types.AttributeValue{
			"PK":            &types.AttributeValueMemberS{Value: pk},
			"SK":            &types.AttributeValueMemberS{Value: sk},
			"Type":          &types.AttributeValueMemberS{Value: "SquarePayment"},
			"TenantID":      &types.AttributeValueMemberS{Value: tenantID},
			"SquareOrderID": &types.AttributeValueMemberS{Value: payment.SquareOrderID},
		},
	}

	_, err := repo.client.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
