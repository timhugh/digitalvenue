package sqs

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	squarePaymentCreatedQueueURL = "SQUARE_PAYMENT_CREATED_QUEUE_URL"
	orderCreatedQueueURL         = "ORDER_CREATED_QUEUE_URL"
)

func NewClient(config aws.Config) *sqs.Client {
	return sqs.NewFromConfig(config)
}
