#!/usr/bin/env bash

if [ -z "${CORE_DATA_TABLE_NAME}" ]; then
  echo "CORE_DATA_TABLE_NAME is not set. Please set the CORE_DATA_TABLE_NAME environment variable."
  exit 1
fi

if [ -z "${SQUARE_MERCHANT_ID}" ]; then
  echo "SQUARE_MERCHANT_ID is not set. Please set the SQUARE_MERCHANT_ID environment variable."
  exit 1
fi

if [ -z "${SQUARE_WEBHOOK_SIGNATURE_KEY}" ]; then
  echo "SQUARE_WEBHOOK_SIGNATURE_KEY is not set. Please set the SQUARE_WEBHOOK_SIGNATURE_KEY environment variable."
  exit 1
fi

if [ -z "${SQUARE_API_ACCESS_TOKEN}" ]; then
  echo "SQUARE_API_ACCESS_TOKEN is not set. Please set the SQUARE_API_ACCESS_TOKEN environment variable."
  exit 1
fi

AWS_PROFILE=${AWS_PROFILE:-default}
echo "Using aws profile '${AWS_PROFILE}'"

aws dynamodb put-item \
  --profile "${AWS_PROFILE}" \
  --table-name "${CORE_DATA_TABLE_NAME}" \
  --item "$(cat <<EOF
{
  "PK":             { "S": "Tenant#tim" },
  "SK":             { "S": "Tenant#tim" },
  "Type":           { "S": "Tenant" },
  "Name":           { "S": "Tim's Test Tenant" },

  "EmailsEnabled":  { "BOOL": true },
  "SMTPAccount":    { "S": "tim@digital-venue.net" },
  "SMTPPassword":   { "S": "password" },
  "SMTPHost":       { "S": "mailcatcher" },
  "SMTPPort":       { "N": "1025" }
}
EOF
)"

aws dynamodb put-item \
  --profile "${AWS_PROFILE}" \
  --table-name "${CORE_DATA_TABLE_NAME}" \
  --item "$(cat <<EOF
{
  "PK":                         { "S": "SquareMerchant#${SQUARE_MERCHANT_ID}" },
  "SK":                         { "S": "SquareMerchant#${SQUARE_MERCHANT_ID}" },
  "Type":                       { "S": "SquareMerchant" },
  "TenantID":                   { "S": "Tenant#tim" },
  "Name":                       { "S": "Test Merchant" },
  "SquareAPIToken":             { "S": "${SQUARE_API_ACCESS_TOKEN}" },
  "SquareWebhookSignatureKey":  { "S": "${SQUARE_WEBHOOK_SIGNATURE_KEY}" }
}
EOF
)"
