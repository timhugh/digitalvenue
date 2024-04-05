#!/usr/bin/env bash

if [ -z "$SQUARE_MERCHANTS_TABLE_NAME" ]; then
  echo "SQUARE_MERCHANTS_TABLE_NAME is not set. Please set the SQUARE_MERCHANTS_TABLE_NAME environment variable."
  exit 1
fi

if [ -z "$SQUARE_PAYMENTS_TABLE_NAME" ]; then
  echo "SQUARE_PAYMENTS_TABLE_NAME is not set. Please set the SQUARE_PAYMENTS_TABLE_NAME environment variable."
  exit 1
fi

if [ -z "$ORDERS_TABLE_NAME" ]; then
  echo "ORDERS_TABLE_NAME is not set. Please set the ORDERS_TABLE_NAME environment variable."
  exit 1
fi

if [ -z "$CUSTOMERS_TABLE_NAME" ]; then
  echo "CUSTOMERS_TABLE_NAME is not set. Please set the CUSTOMERS_TABLE_NAME environment variable."
  exit 1
fi

aws dynamodb create-table --table-name "$SQUARE_MERCHANTS_TABLE_NAME" --attribute-definitions AttributeName=SquareMerchantID,AttributeType=S  --key-schema AttributeName=SquareMerchantID,KeyType=HASH  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://localhost:8000
aws dynamodb create-table --table-name "$SQUARE_PAYMENTS_TABLE_NAME"  --attribute-definitions AttributeName=SquarePaymentID,AttributeType=S   --key-schema AttributeName=SquarePaymentID,KeyType=HASH   --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://localhost:8000
aws dynamodb create-table --table-name "$ORDERS_TABLE_NAME"           --attribute-definitions AttributeName=OrderID,AttributeType=S           --key-schema AttributeName=OrderID,KeyType=HASH           --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://localhost:8000
aws dynamodb create-table --table-name "$CUSTOMERS_TABLE_NAME"        --attribute-definitions AttributeName=CustomerID,AttributeType=S        --key-schema AttributeName=CustomerID,KeyType=HASH        --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://localhost:8000
