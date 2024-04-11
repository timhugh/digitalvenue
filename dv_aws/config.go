package dv_aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/timhugh/digitalvenue/core"
)

func DefaultConfig() (aws.Config, error) {
	localURL, err := core.RequireEnv("LOCAL_DYNAMODB_URL")
	if err == nil {
		return LocalConfig(localURL)
	}
	return config.LoadDefaultConfig(context.TODO())
}

func LocalConfig(localURL string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: localURL}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
}
