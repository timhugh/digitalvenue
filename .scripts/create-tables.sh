#!/usr/bin/env bash

if [ -z "$CORE_DATA_TABLE_NAME" ]; then
  echo "CORE_DATA_TABLE_NAME is not set. Please set the CORE_DATA_TABLE_NAME environment variable."
  exit 1
fi

AWS_PROFILE=${AWS_PROFILE:-default}
echo "Using aws profile '${AWS_PROFILE}' (Use AWS_PROFILE to change)"
read -p "Continue? (y/N)> " CONT
if [ "$CONT" != "y" ]; then
  echo "Exiting"
  exit 0
fi

aws dynamodb create-table \
  --profile ${AWS_PROFILE} \
  --table-name "${CORE_DATA_TABLE_NAME}" \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=CustomerID,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --local-secondary-indexes '[
    {
      "IndexName": "CustomerIDIndex",
      "KeySchema": [
        {
          "AttributeName": "PK",
          "KeyType": "HASH"
        },
        {
          "AttributeName": "CustomerID",
          "KeyType": "RANGE"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      }
    }
  ]'
