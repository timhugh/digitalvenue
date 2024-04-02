package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/db"
	"os"
)

type MerchantsRepositoryConfig struct {
	TableName string
}

func NewMerchantsRepositoryConfig() MerchantsRepositoryConfig {
	return MerchantsRepositoryConfig{
		TableName: os.Getenv(MerchantsTableName),
	}
}

type MerchantsRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewMerchantsRepository(config MerchantsRepositoryConfig, client *dynamodb.Client) db.MerchantsRepository {
	return MerchantsRepository{
		tableName: config.TableName,
		client:    client,
	}
}

func (repo MerchantsRepository) CreateMerchant(merchant core.Merchant) error {
	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			SquareMerchantId:          &types.AttributeValueMemberS{Value: merchant.SquareMerchantId},
			SquareWebhookSignatureKey: &types.AttributeValueMemberS{Value: merchant.SquareWebhookSignatureKey},
			SquareAPIKey:              &types.AttributeValueMemberS{Value: merchant.SquareAPIKey},
		},
		TableName: aws.String(repo.tableName),
	}
	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return fmt.Errorf("failed to create merchant: %w", err)
	}

	return nil
}

func (repo MerchantsRepository) FindMerchantBySquareMerchantId(squareMerchantId string) (core.Merchant, error) {
	var merchant = core.Merchant{}

	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			SquareMerchantId: &types.AttributeValueMemberS{Value: squareMerchantId},
		},
		TableName: aws.String(repo.tableName),
	}

	getItemOutput, err := repo.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return merchant, fmt.Errorf("unable to retrieve merchant with id '%s': %w", squareMerchantId, err)
	}

	err = attributevalue.UnmarshalMap(getItemOutput.Item, &merchant)

	return merchant, err
}
