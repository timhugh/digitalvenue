#!/usr/bin/env bash

if [ -z "$MERCHANTS_TABLE_NAME" ]; then
  echo "MERCHANTS_TABLE_NAME is not set. Please set the MERCHANTS_TABLE_NAME environment variable."
  exit 1
fi

if [ -z "$PAYMENTS_TABLE_NAME" ]; then
  echo "PAYMENTS_TABLE_NAME is not set. Please set the PAYMENTS_TABLE_NAME environment variable."
  exit 1
fi

aws dynamodb create-table --table-name "$MERCHANTS_TABLE_NAME" --attribute-definitions AttributeName=SquareMerchantId,AttributeType=S --key-schema AttributeName=SquareMerchantId,KeyType=HASH --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://localhost:8000
aws dynamodb create-table --table-name "$PAYMENTS_TABLE_NAME"  --attribute-definitions AttributeName=SquarePaymentId,AttributeType=S  --key-schema AttributeName=SquarePaymentId,KeyType=HASH  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://localhost:8000
