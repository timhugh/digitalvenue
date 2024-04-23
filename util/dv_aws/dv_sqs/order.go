package dv_sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
)

type OrderProcessedQueue struct {
	sqsClient Client
	queueURL  string
}

func NewOrderProcessedQueue(sqsClient Client) (*OrderProcessedQueue, error) {
	queueURL, err := core.RequireEnv("ORDER_PROCESSED_QUEUE_URL")
	if err != nil {
		return nil, err
	}

	return &OrderProcessedQueue{sqsClient, queueURL}, nil
}

func (q *OrderProcessedQueue) PublishOrderProcessedEvent(tenantID string, orderID string) error {
	payload := tenantID + ":" + orderID
	_, err := q.sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: &payload,
	})
	if err != nil {
		return errors.Wrap(err, "failed to publish OrderProcessedEvent")
	}

	return nil
}
