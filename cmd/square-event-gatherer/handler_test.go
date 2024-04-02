package main

import (
	"github.com/matryer/is"
	"github.com/rs/zerolog"
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

			handler := newHandler(log)

			response, err := handler.handle(testCase.request)
			is.NoErr(err)
			is.Equal(testCase.expectedResponse, response)
		})
	}
}
