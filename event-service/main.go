package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting events-service")

	handler, err := initializeHandler()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create handler")
	}

	defer log.Info().Msg("Exiting events-service")
	lambda.Start(handler.handle)
}
