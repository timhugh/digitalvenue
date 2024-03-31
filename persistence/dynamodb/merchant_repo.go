package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
)

const (
	squareMerchantId          = "SquareMerchantId"
	squareWebhookSignatureKey = "SquareWebhookSignatureKey"
	squareAPIKey              = "SquareAPIKey"
)

type MerchantRepo struct {
	client    *dynamodb.Client
	tableName string
}

func NewMerchantRepo(tableName string) (*MerchantRepo, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	//awsConfig, err := config.LoadDefaultConfig(context.TODO(),
	//	config.WithRegion("us-west-2"),
	//	config.WithEndpointResolver(aws.EndpointResolverFunc(
	//		func(service, region string) (aws.Endpoint, error) {
	//			return aws.Endpoint{URL: "http://localhost:8000"}, nil
	//		})),
	//	config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
	//		Value: aws.Credentials{
	//			AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
	//			Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
	//		},
	//	}),
	//)
	if err != nil {
		return nil, fmt.Errorf("failed to configure aws: %w", err)
	}

	return &MerchantRepo{
		client:    dynamodb.NewFromConfig(awsConfig),
		tableName: tableName,
	}, nil
}

func (r *MerchantRepo) CreateMerchant(merchant *core.Merchant) error {
	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			squareMerchantId:          &types.AttributeValueMemberS{Value: merchant.SquareMerchantId},
			squareWebhookSignatureKey: &types.AttributeValueMemberS{Value: merchant.SquareWebhookSignatureKey},
			squareAPIKey:              &types.AttributeValueMemberS{Value: merchant.SquareAPIKey},
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

func (r *MerchantRepo) FindMerchantBySquareMerchantId(squareMerchantId string) (*core.Merchant, error) {
	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			squareMerchantId: &types.AttributeValueMemberS{Value: squareMerchantId},
		},
		TableName: aws.String(r.tableName),
	}
	getItemResponse, err := r.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get merchant")
		return nil, fmt.Errorf("unable to retrieve merchant with id '%s'", squareMerchantId)
	}

	var merchant *core.Merchant
	err = attributevalue.UnmarshalMap(getItemResponse.Item, merchant)

	return merchant, nil
}
