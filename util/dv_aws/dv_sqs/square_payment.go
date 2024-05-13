package dv_sqs

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/square"
)

type SquarePaymentCreatedEvent struct {
	SquarePaymentID  string `json:"square_payment_id"`
	SquareMerchantID string `json:"square_merchant_id"`
}

type SquarePaymentCreatedQueue struct {
	*Queue
}

func NewSquarePaymentCreatedQueue(client Client) (*SquarePaymentCreatedQueue, error) {
	queueName, err := core.RequireEnv("SQUARE_PAYMENT_CREATED_EVENT_QUEUE_URL")
	if err != nil {
		return nil, err
	}

	return &SquarePaymentCreatedQueue{NewQueue(client, queueName)}, nil
}

func (q *SquarePaymentCreatedQueue) PublishPaymentCreated(payment *square.Payment) error {
	event := SquarePaymentCreatedEvent{
		SquarePaymentID:  payment.SquarePaymentID,
		SquareMerchantID: payment.SquareMerchantID,
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal SquarePaymentCreatedEvent")
	}

	_, err = q.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: aws.String(string(eventJSON)),
	})
	if err != nil {
		return errors.Wrap(err, "failed to publish SquarePaymentCreatedEvent")
	}

	return nil
}
