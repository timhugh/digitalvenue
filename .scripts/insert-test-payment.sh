#!/usr/bin/env bash

if [ -z "$SQUARE_PAYMENTS_TABLE_NAME" ]; then
  echo "$SQUARE_PAYMENTS_TABLE_NAME is not set. Please set the $SQUARE_PAYMENTS_TABLE_NAME environment variable."
  exit 1
fi

if [ -z "$SQUARE_PAYMENT_ID" ]; then
  echo "$SQUARE_PAYMENT_ID is not set. Please set the $SQUARE_PAYMENT_ID environment variable."
  exit 1
fi

if [ -z "$SQUARE_MERCHANT_ID" ]; then
  echo "$SQUARE_MERCHANT_ID is not set. Please set the $SQUARE_MERCHANT_ID environment variable."
  exit 1
fi

if [ -z "$SQUARE_ORDER_ID" ]; then
  echo "$SQUARE_ORDER_ID is not set. Please set the $SQUARE_ORDER_ID environment variable."
  exit 1
fi

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name "$SQUARE_PAYMENTS_TABLE_NAME" --item "$(cat <<EOF
{
  "PK": {"S": "SquareMerchant#$SQUARE_MERCHANT_ID"},
  "SK": {"S": "SquarePayment#$SQUARE_PAYMENT_ID"},
  "Type": {"S": "SquarePayment"},
  "TenantID": {"S": "test"},
  "SquareOrderID": {"S": "$SQUARE_ORDER_ID"}
}
EOF
)"
