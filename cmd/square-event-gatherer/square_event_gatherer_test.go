package main

import (
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/logger"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func buildSuccessRecord() events.DynamoDBEventRecord {
	return events.DynamoDBEventRecord{
		EventID:   "event-id",
		EventName: "INSERT",
		Change: events.DynamoDBStreamRecord{
			NewImage: map[string]events.DynamoDBAttributeValue{
				"Type": events.NewStringAttribute("SquarePayment"),
				"PK":   events.NewStringAttribute("Merchant#" + squaretest.SquareMerchantID),
				"SK":   events.NewStringAttribute("SquarePayment#" + squaretest.SquarePaymentID),
			},
		},
	}
}

func buildWrongEventRecord() events.DynamoDBEventRecord {
	return events.DynamoDBEventRecord{
		EventID:   "event-id",
		EventName: "WRONG",
	}
}

func buildWrongTypeRecord() events.DynamoDBEventRecord {
	return events.DynamoDBEventRecord{
		EventID:   "event-id",
		EventName: "INSERT",
		Change: events.DynamoDBStreamRecord{
			NewImage: map[string]events.DynamoDBAttributeValue{
				"Type": events.NewStringAttribute("NotSquarePayment"),
			},
		},
	}
}

func TestSquareEventGathererHandler(t *testing.T) {
	testCases := []struct {
		name             string
		request          events.DynamoDBEvent
		expectedResponse events.DynamoDBEventResponse
	}{
		{
			name: "basic success",
			request: events.DynamoDBEvent{
				Records: []events.DynamoDBEventRecord{
					buildSuccessRecord(),
				},
			},
			expectedResponse: events.DynamoDBEventResponse{},
		},
		{
			name: "wrong event",
			request: events.DynamoDBEvent{
				Records: []events.DynamoDBEventRecord{
					buildWrongEventRecord(),
				},
			},
			expectedResponse: events.DynamoDBEventResponse{},
		},
		{
			name: "wrong type",
			request: events.DynamoDBEvent{
				Records: []events.DynamoDBEventRecord{
					buildWrongTypeRecord(),
				},
			},
			expectedResponse: events.DynamoDBEventResponse{},
		},
	}

	log := logger.Default()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			is := is.New(t)
			mock.SetUp(t)

			gatherer := mock.Mock[square.PaymentGatherer]()
			handler := NewSquareEventGathererHandler(log, gatherer)

			response, err := handler.Handle(testCase.request)
			is.NoErr(err)
			is.Equal(testCase.expectedResponse, response)
		})
	}
}
