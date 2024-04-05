package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/square"
)

func (repo *Repository) PutSquareMerchant(merchant square.Merchant) error {
	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			SquareMerchantID:          &types.AttributeValueMemberS{Value: merchant.SquareMerchantID},
			SquareWebhookSignatureKey: &types.AttributeValueMemberS{Value: merchant.SquareWebhookSignatureKey},
			SquareAPIToken:            &types.AttributeValueMemberS{Value: merchant.SquareAPIToken},
		},
		TableName: aws.String(repo.squareMerchantsTableName),
	}
	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetSquareMerchant(squareMerchantID string) (square.Merchant, error) {
	var merchant = square.Merchant{}

	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			SquareMerchantID: &types.AttributeValueMemberS{Value: squareMerchantID},
		},
		TableName: aws.String(repo.squareMerchantsTableName),
	}

	getItemOutput, err := repo.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return merchant, fmt.Errorf("unable to retrieve merchant with id '%s': %w", squareMerchantID, err)
	}

	err = attributevalue.UnmarshalMap(getItemOutput.Item, &merchant)
	if err != nil {
		return merchant, fmt.Errorf("failed to unmarshal merchant: %w", err)
	}

	return merchant, err
}
