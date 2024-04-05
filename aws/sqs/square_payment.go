package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/timhugh/digitalvenue/core"
)

type SquarePaymentCreatedQueue struct {
	client   *sqs.Client
	queueURL string
}

func NewSquarePaymentCreatedQueue(client *sqs.Client) *SquarePaymentCreatedQueue {
	return &SquarePaymentCreatedQueue{
		client:   client,
		queueURL: core.Getenv(squarePaymentCreatedQueueURL),
	}
}

func (queue *SquarePaymentCreatedQueue) PublishSquarePaymentCreated(squarePaymentID string) error {
	_, err := queue.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(squarePaymentID),
		QueueUrl:    aws.String(queue.queueURL),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to sqs: %w", err)
	}

	return nil
}
