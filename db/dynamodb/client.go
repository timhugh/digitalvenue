package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var client *dynamodb.Client

func Connect() (*dynamodb.Client, error) {
	if client == nil {
		cfg, err := DefaultConfig()
		if err != nil {
			return nil, err
		}
		client = dynamodb.NewFromConfig(cfg)
	}
	return client, nil
}

func DefaultConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO())
}

// For testing with local DynamoDB
func LocalConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
}