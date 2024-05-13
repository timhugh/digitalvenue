package dv_sqs

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
)

type OrderProcessedEvent struct {
	OrderID  string `json:"order_id"`
	TenantID string `json:"tenant_id"`
}

type OrderProcessedQueue struct {
	*Queue
}

func NewOrderProcessedQueue(sqsClient Client) (*OrderProcessedQueue, error) {
	queueURL, err := core.RequireEnv("ORDER_PROCESSED_QUEUE_URL")
	if err != nil {
		return nil, err
	}

	return &OrderProcessedQueue{
		NewQueue(sqsClient, queueURL),
	}, nil
}

func (q *OrderProcessedQueue) PublishOrderProcessedEvent(tenantID string, orderID string) error {
	event := OrderProcessedEvent{
		OrderID:  orderID,
		TenantID: tenantID,
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal OrderProcessedEvent")
	}

	_, err = q.sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: aws.String(string(eventJSON)),
	})
	if err != nil {
		return errors.Wrap(err, "failed to publish OrderProcessedEvent")
	}

	return nil
}

type OrderCreatedEvent struct {
	OrderID  string `json:"order_id"`
	TenantID string `json:"tenant_id"`
}

type OrderCreatedQueue struct {
	*Queue
}

func NewOrderCreatedQueue(sqsClient Client) (*OrderCreatedQueue, error) {
	queueURL, err := core.RequireEnv("ORDER_CREATED_QUEUE_URL")
	if err != nil {
		return nil, err
	}

	return &OrderCreatedQueue{
		NewQueue(sqsClient, queueURL),
	}, nil
}

func (q *OrderCreatedQueue) PublishOrderCreatedEvent(tenantID string, orderID string) error {
	event := OrderCreatedEvent{
		OrderID:  orderID,
		TenantID: tenantID,
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal OrderCreatedEvent")
	}

	_, err = q.sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: aws.String(string(eventJSON)),
	})
	if err != nil {
		return errors.Wrap(err, "failed to publish OrderCreatedEvent")
	}

	return nil
}
