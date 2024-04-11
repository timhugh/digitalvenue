package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/square"
)

type squareMerchant struct {
	PK                        string
	SK                        string
	Type                      string
	TenantID                  string
	Name                      string
	SquareAPIToken            string
	SquareWebhookSignatureKey string
}

func (repo *Repository) GetSquareMerchant(squareMerchantID string) (*square.Merchant, error) {
	merchantKey := "SquareMerchant#" + squareMerchantID
	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: merchantKey},
			"SK": &types.AttributeValueMemberS{Value: merchantKey},
		},
	}

	item := squareMerchant{}
	err := repo.getItem("SquareMerchant", input, &item)
	if err != nil {
		return nil, err
	}

	tenantID, err := removeIDPrefix(item.TenantID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid tenant ID")
	}

	return &square.Merchant{
		TenantID:                  tenantID,
		Name:                      item.Name,
		ID:                        squareMerchantID,
		SquareWebhookSignatureKey: item.SquareWebhookSignatureKey,
		SquareAPIToken:            item.SquareAPIToken,
	}, nil
}
