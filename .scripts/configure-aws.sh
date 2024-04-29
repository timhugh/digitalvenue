#!/usr/bin/env bash

if [ -z "${LOCAL_AWS_URL}" ]; then
  echo "LOCAL_AWS_URL is not set. Please set the LOCAL_AWS_URL environment variable."
  exit 1
fi

aws --profile=localstack configure set aws_access_key_id      test
aws --profile=localstack configure set aws_secret_access_key  test
aws --profile=localstack configure set region                 us-west-2
aws --profile=localstack configure set endpoint_url           "${LOCAL_AWS_URL}"
