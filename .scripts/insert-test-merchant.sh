#!/usr/bin/env bash

if [ -z "$MERCHANTS_TABLE" ]; then
  echo "MERCHANTS_TABLE is not set. Please set the MERCHANTS_TABLE environment variable."
  exit 1
fi

if [ -z "$SQUARE_MERCHANT_ID" ]; then
  echo "SQUARE_MERCHANT_ID is not set. Please set the SQUARE_MERCHANT_ID environment variable."
  exit 1
fi

if [ -z "$SQUARE_WEBHOOK_SIGNATURE_KEY" ]; then
  echo "SQUARE_WEBHOOK_SIGNATURE_KEY is not set. Please set the SQUARE_WEBHOOK_SIGNATURE_KEY environment variable."
  exit 1
fi

if [ -z "$SQUARE_API_ACCESS_TOKEN" ]; then
  echo "SQUARE_API_ACCESS_TOKEN is not set. Please set the SQUARE_API_ACCESS_TOKEN environment variable."
  exit 1
fi

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name "$MERCHANTS_TABLE" --item "$(cat <<EOF
{
  "SquareMerchantId": {"S": "$SQUARE_MERCHANT_ID"},
  "SquareWebhookSignatureKey": {"S": "$SQUARE_WEBHOOK_SIGNATURE_KEY"},
  "SquareAPIKey": {"S": "$SQUARE_API_ACCESS_TOKEN"}
}
EOF
)"