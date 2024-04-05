package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"os"
)

type Queue struct {
	client *sqs.Client

	squarePaymentCreatedQueueURL string
}

func NewQueue(client *sqs.Client) *Queue {
	return &Queue{
		client:                       client,
		squarePaymentCreatedQueueURL: os.Getenv(squarePaymentCreatedQueueURL),
	}
}

func (queue *Queue) PublishSquarePaymentCreated(squarePaymentID string) error {
	_, err := queue.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(squarePaymentID),
		QueueUrl:    aws.String(queue.squarePaymentCreatedQueueURL),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to sqs: %w", err)
	}

	return nil
}
