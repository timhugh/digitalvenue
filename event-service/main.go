package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/db/dynamodb"
	"os"
)

func main() {
	log.Info().Msg("Starting events-service")

	config := NewEventServiceConfig()
	merchantRepo, err := dynamodb.NewMerchantsRespository(os.Getenv("MERCHANTS_TABLE"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create merchant repository")
	}

	handler := handler(config, merchantRepo)
	lambda.Start(handler)
}
