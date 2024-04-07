package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/core"
)

type SquarePaymentCreatedQueue struct {
	client   *sqs.Client
	queueURL string
}

func NewSquarePaymentCreatedQueue(client *sqs.Client) (*SquarePaymentCreatedQueue, error) {
	tableName, err := core.RequireEnv(squarePaymentCreatedQueueURL)
	if err != nil {
		return nil, err
	}

	return &SquarePaymentCreatedQueue{
		client:   client,
		queueURL: tableName,
	}, nil
}

func (queue *SquarePaymentCreatedQueue) PublishSquarePaymentCreated(squarePaymentID string) error {
	_, err := queue.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(squarePaymentID),
		QueueUrl:    aws.String(queue.queueURL),
	})
	if err != nil {
		return errors.Wrap(err, "failed to send message to sqs")
	}

	return nil
}
