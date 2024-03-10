#!/usr/bin/env sh

set -x

aws cloudformation deploy --template-file .cloudformation/lambda-role.yml --stack-name digitalvenue-lambda-role --capabilities CAPABILITY_IAM
aws cloudformation deploy --template-file .cloudformation/lambda-gateway.yml --stack-name digitalvenue-lambda-gateway

# lambdaTest
aws cloudformation deploy --template-file lambda-test/ecr.yml --stack-name digitalvenue-lambda-test-ecr
docker build -t 392387111634.dkr.ecr.us-west-2.amazonaws.com/digitalvenue/lambda-test:latest lambda-test
docker push 392387111634.dkr.ecr.us-west-2.amazonaws.com/digitalvenue/lambda-test:latest
aws cloudformation deploy --template-file lambda-test/lambda.yml --stack-name digitalvenue-lambda-test \
  --parameter-overrides LambdaTestECRStackname=digitalvenue-lambda-test-ecr LambdaRoleStackname=digitalvenue-lambda-role
aws cloudformation deploy --template-file lambda-test/api.yml --stack-name digitalvenue-lambda-test-api \
  --parameter-overrides LambdaGatewayStackName=digitalvenue-lambda-gateway LambdaTestStackName=digitalvenue-lambda-test

aws cloudformation deploy --template-file .cloudformation/lambda-gateway-deployment.yml --stack-name digitalvenue-lambda-gateway-deployment \
  --parameter-overrides LambdaGatewayStackName=digitalvenue-lambda-gateway



# Cleanup:
#aws cloudformation delete-stack --stack-name digitalvenue-lambda-gateway-deployment
#aws cloudformation delete-stack --stack-name digitalvenue-lambda-test-api
#aws cloudformation delete-stack --stack-name digitalvenue-lambda-test
#aws cloudformation delete-stack --stack-name digitalvenue-lambda-test-ecr
#aws cloudformation delete-stack --stack-name digitalvenue-lambda-gateway
#aws cloudformation delete-stack --stack-name digitalvenue-lambda-role
