package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/core"
)

type OrderCreatedQueue struct {
	client   *sqs.Client
	queueURL string
}

func NewOrderCreatedQueue(client *sqs.Client) (*OrderCreatedQueue, error) {
	tableName, err := core.RequireEnv(orderCreatedQueueURL)
	if err != nil {
		return nil, err
	}

	return &OrderCreatedQueue{
		client:   client,
		queueURL: tableName,
	}, nil
}

func (queue *OrderCreatedQueue) Publish(orderID string) error {
	_, err := queue.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(orderID),
		QueueUrl:    aws.String(queue.queueURL),
	})
	if err != nil {
		return errors.Wrap(err, "failed to send message to sqs")
	}

	return nil
}
