package main

import (
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name             string
		request          events.SQSEvent
		expectedResponse events.SQSEventResponse
	}{
		{
			name: "basic success",
			request: events.SQSEvent{
				Records: []events.SQSMessage{
					{
						Body: `{"type": "payment.created"}`,
					},
				},
			},
		},
	}

	log := zerolog.Logger{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			is := is.New(t)
			mock.SetUp(t)

			gatherer := mock.Mock[square.EventGatherer]()
			handler := newHandler(log, gatherer)

			response, err := handler.handle(testCase.request)
			is.NoErr(err)
			is.Equal(testCase.expectedResponse, response)
		})
	}
}
