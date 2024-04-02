package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
	"os"
)

const (
	SquareMerchantId          = "SquareMerchantId"
	SquareWebhookSignatureKey = "SquareWebhookSignatureKey"
	SquareAPIKey              = "SquareAPIKey"
)

type MerchantsRepositoryConfig struct {
	TableName string
}

func NewMerchantsRepositoryConfig() MerchantsRepositoryConfig {
	return MerchantsRepositoryConfig{
		TableName: os.Getenv("MERCHANTS_TABLE"),
	}
}

type MerchantsRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewMerchantsRepository(config MerchantsRepositoryConfig, client *dynamodb.Client) MerchantsRepository {
	return MerchantsRepository{
		tableName: config.TableName,
		client:    client,
	}
}

func (r MerchantsRepository) CreateMerchant(merchant core.Merchant) error {
	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			SquareMerchantId:          &types.AttributeValueMemberS{Value: merchant.SquareMerchantId},
			SquareWebhookSignatureKey: &types.AttributeValueMemberS{Value: merchant.SquareWebhookSignatureKey},
			SquareAPIKey:              &types.AttributeValueMemberS{Value: merchant.SquareAPIKey},
		},
		TableName: aws.String(r.tableName),
	}
	_, err := r.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create merchant")
		return fmt.Errorf("failed to create merchant")
	}

	return nil
}

func (r MerchantsRepository) FindMerchantBySquareMerchantId(squareMerchantId string) (core.Merchant, error) {
	var merchant = core.Merchant{}

	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			SquareMerchantId: &types.AttributeValueMemberS{Value: squareMerchantId},
		},
		TableName: aws.String(r.tableName),
	}

	getItemOutput, err := r.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get merchant")
		return merchant, fmt.Errorf("unable to retrieve merchant with id '%s'", squareMerchantId)
	}

	err = attributevalue.UnmarshalMap(getItemOutput.Item, &merchant)

	return merchant, err
}
