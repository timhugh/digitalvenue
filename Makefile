include .env

app_name = digitalvenue
environment = dev
code_bucket = $(app_name)-codebucket

.PHONY: validate
validate:
	aws cloudformation validate-template --template-body file://./.cloudformation/env.yml

.PHONY: codebucket
codebucket:
	aws cloudformation deploy --stack-name $(code_bucket) \
		--template-file ./.cloudformation/bucket.yml \
		--parameter-overrides \
			BucketName=$(code_bucket)

.PHONY: package
package: codebucket build
	aws cloudformation package \
		--template-file ./.cloudformation/env.yml \
		--output-template-file ./template.yml \
		--s3-bucket $(code_bucket) \
		--s3-prefix $(app_name)/$(environment)

.PHONY: deploy
deploy: package
	aws cloudformation deploy \
		--stack-name $(app_name)-$(environment) \
		--template-file ./template.yml \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			Route53HostedZoneId=${ROUTE_53_HOSTED_ZONE_ID} \
			Environment=$(environment) \
			CodeBucketName=$(code_bucket)

.PHONY: build
build: build/hello-world.zip

build/%.zip: %
	GOOS=linux GOARCH=arm64 go build -o ../build/$^/bootstrap -C $^
	zip -rj $@ build/$^/bootstrap

clean:
	rm -f template.yml
	rm -rf build/
