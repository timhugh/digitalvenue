#!/usr/bin/env bash

if [ -z "$SQUARE_MERCHANTS_TABLE_NAME" ]; then
  echo "$SQUARE_MERCHANTS_TABLE_NAME is not set. Please set the $SQUARE_MERCHANTS_TABLE_NAME environment variable."
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

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name "$SQUARE_MERCHANTS_TABLE_NAME" --item "$(cat <<EOF
{
  "TenantID": {"S": "test"},
  "SquareMerchantID": {"S": "$SQUARE_MERCHANT_ID"},
  "SquareWebhookSignatureKey": {"S": "$SQUARE_WEBHOOK_SIGNATURE_KEY"},
  "SquareAPIToken": {"S": "$SQUARE_API_ACCESS_TOKEN"}
}
EOF
)"
