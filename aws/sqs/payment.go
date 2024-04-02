package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/timhugh/digitalvenue/queue"
	"os"
)

type PaymentCreatedQueueConfig struct {
	QueueUrl string
}

func NewPaymentCreatedQueueConfig() PaymentCreatedQueueConfig {
	return PaymentCreatedQueueConfig{
		QueueUrl: os.Getenv("PAYMENT_CREATED_QUEUE_URL"),
	}
}

type PaymentCreatedQueue struct {
	queueUrl string
	client   *sqs.Client
}

func NewPaymentCreatedQueue(config PaymentCreatedQueueConfig, client *sqs.Client) queue.PaymentCreatedQueue {
	return PaymentCreatedQueue{
		queueUrl: config.QueueUrl,
		client:   client,
	}
}

func (queue PaymentCreatedQueue) Publish(squarePaymentId string) error {
	_, err := queue.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(string(squarePaymentId)),
		QueueUrl:    aws.String(queue.queueUrl),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to sqs: %w", err)
	}

	return nil
}
