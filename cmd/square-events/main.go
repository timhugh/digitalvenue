package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

func main() {
	handler, err := initializeHandler()
	if err != nil {
		log.Fatal().Err(err).Str("service", "square-events").Msg("Failed to initialize handler")
	}
	lambda.Start(handler.Handle)
}
