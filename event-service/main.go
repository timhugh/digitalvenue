package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler, err := initializeHandler()
	if err != nil {
		panic("Failed to initialize handler: " + err.Error())
	}
	lambda.Start(handler.handle)
}
