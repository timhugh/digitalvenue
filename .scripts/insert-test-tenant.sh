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

SMTP_ACCOUNT="${SMTP_ACCOUNT:-"tim@digital-venue.net"}"
SMTP_PASSWORD="${SMTP_PASSWORD:-"password"}"
SMTP_HOST="${SMTP_HOST:-"mailcatcher"}"
SMTP_PORT="${SMTP_PORT:-"1025"}"

AWS_PROFILE=${AWS_PROFILE:-default}
echo "Using aws profile '${AWS_PROFILE}' (Use AWS_PROFILE to change)"
read -p "Continue? (y/N)> " CONT
if [ "$CONT" != "y" ]; then
  echo "Exiting"
  exit 0
fi

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
  "SMTPAccount":    { "S": "${SMTP_ACCOUNT}" },
  "SMTPPassword":   { "S": "${SMTP_PASSWORD}" },
  "SMTPHost":       { "S": "${SMTP_HOST}" },
  "SMTPPort":       { "N": "${SMTP_PORT}" },
  "SMTPFromAddress":{ "S": "${SMTP_FROM_ADDRESS}" }
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
  "SquareWebhookSignatureKey":  { "S": "${SQUARE_WEBHOOK_SIGNATURE_KEY}" },
  "TicketableCategories":       { "L": [] }
}
EOF
)"

aws s3 cp lcgf/ticketEmail.html s3://${S3_TEMPLATE_BUCKET_NAME}/tim/
aws s3 cp lcgf/facebook.png s3://${S3_TENANT_FILES_BUCKET_NAME}/tim/ --acl public-read
aws s3 cp lcgf/instagram.png s3://${S3_TENANT_FILES_BUCKET_NAME}/tim/ --acl public-read
aws s3 cp lcgf/ticket-header.jpg s3://${S3_TENANT_FILES_BUCKET_NAME}/tim/ --acl public-read
