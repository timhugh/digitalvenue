//go:build !localDynamoDB
// +build !localDynamoDB

package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO())
}
