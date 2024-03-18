#!/usr/bin/env sh

docker buildx build --push --provenance false --platform linux/arm64 \
  -t ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com/digitalvenue/runtime-base:latest \
  --file docker/builder-base.Dockerfile .
