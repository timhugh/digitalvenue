package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/timhugh/digitalvenue/square/queue"
	"os"
)

type PaymentCreatedQueueConfig struct {
	QueueURL string
}

func NewPaymentCreatedQueueConfig() PaymentCreatedQueueConfig {
	return PaymentCreatedQueueConfig{
		QueueURL: os.Getenv("PAYMENT_CREATED_QUEUE_URL"),
	}
}

type PaymentCreatedQueue struct {
	queueURL string
	client   *sqs.Client
}

func NewPaymentCreatedQueue(config PaymentCreatedQueueConfig, client *sqs.Client) queue.PaymentCreatedQueue {
	return PaymentCreatedQueue{
		queueURL: config.QueueURL,
		client:   client,
	}
}

func (queue PaymentCreatedQueue) Publish(squarePaymentID string) error {
	_, err := queue.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(string(squarePaymentID)),
		QueueUrl:    aws.String(queue.queueURL),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to sqs: %w", err)
	}

	return nil
}
