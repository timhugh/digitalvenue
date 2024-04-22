package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"os"
)

func main() {
	log := logger.Default().AddParam("service", "square-event-gatherer")
	env, err := core.RequireEnv("ENVIRONMENT")
	if err != nil {
		log.AddParam("error", err).Fatal("Failed to determine application environment")
		os.Exit(1)
	}
	log.AddParam("environment", env)

	handler, err := initializeHandler(log)
	if err != nil {
		log.AddParam("error", err).Fatal("Failed to initialize handler")
		os.Exit(1)
	}
	lambda.Start(handler.Handle)
}
