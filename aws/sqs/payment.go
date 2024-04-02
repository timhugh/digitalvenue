package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/timhugh/digitalvenue/square/queue"
	"os"
)

type SquarePaymentCreatedQueueConfig struct {
	QueueURL string
}

func NewSquarePaymentCreatedQueueConfig() SquarePaymentCreatedQueueConfig {
	return SquarePaymentCreatedQueueConfig{
		QueueURL: os.Getenv("SQUARE_PAYMENT_CREATED_QUEUE_URL"),
	}
}

type SquarePaymentCreatedQueue struct {
	queueURL string
	client   *sqs.Client
}

func NewSquarePaymentCreatedQueue(config SquarePaymentCreatedQueueConfig, client *sqs.Client) queue.SquarePaymentCreatedQueue {
	return SquarePaymentCreatedQueue{
		queueURL: config.QueueURL,
		client:   client,
	}
}

func (queue SquarePaymentCreatedQueue) Publish(squarePaymentID string) error {
	_, err := queue.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(string(squarePaymentID)),
		QueueUrl:    aws.String(queue.queueURL),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to sqs: %w", err)
	}

	return nil
}
