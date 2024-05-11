package dv_sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client interface {
	SendMessage(ctx context.Context, input *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

func NewClient(config aws.Config) Client {
	return sqs.NewFromConfig(config)
}

type Queue struct {
	sqsClient Client
	queueURL  string
}

func NewQueue(sqsClient Client, queueURL string) *Queue {
	return &Queue{
		sqsClient: sqsClient,
		queueURL:  queueURL,
	}
}
