#!/usr/bin/env bash

aws --profile=localstack configure set aws_access_key_id      test
aws --profile=localstack configure set aws_secret_access_key  test
aws --profile=localstack configure set region                 us-west-2
aws --profile=localstack configure set endpoint_url           http://localhost:4566
