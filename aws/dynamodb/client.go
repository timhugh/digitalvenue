package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewClient(config aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(config)
}
