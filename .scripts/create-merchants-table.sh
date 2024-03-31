#!/usr/bin/env bash

if [ -z "$MERCHANTS_TABLE" ]; then
  echo "MERCHANTS_TABLE is not set. Please set the MERCHANTS_TABLE environment variable."
  exit 1
fi

aws dynamodb create-table --table-name "$MERCHANTS_TABLE" --attribute-definitions AttributeName=SquareMerchantId,AttributeType=S --key-schema AttributeName=SquareMerchantId,KeyType=HASH --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://localhost:8000
