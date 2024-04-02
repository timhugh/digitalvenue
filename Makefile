include .env

app_name = digitalvenue
environment = dev
code_bucket = $(app_name)-codebucket
root = $(shell git rev-parse --show-toplevel)

.PHONY: validate
validate:
	aws cloudformation validate-template --template-body file://$(root)/.cloudformation/env.yml

.PHONY: codebucket
codebucket:
	aws cloudformation deploy --stack-name $(code_bucket) \
		--template-file $(root)/.cloudformation/bucket.yml \
		--parameter-overrides \
			BucketName=$(code_bucket)

.PHONY: package
package: codebucket build
	aws cloudformation package \
		--template-file $(root)/.cloudformation/env.yml \
		--output-template-file $(root)/template.yml \
		--s3-bucket $(code_bucket) \
		--s3-prefix $(app_name)/$(environment)

.PHONY: deploy
deploy: package
	aws cloudformation deploy \
		--stack-name $(app_name)-$(environment) \
		--template-file $(root)/template.yml \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			Route53HostedZoneId=${ROUTE_53_HOSTED_ZONE_ID} \
			Environment=$(environment) \
			CodeBucketName=$(code_bucket)

.PHONY: build
build: build/echo.zip build/square-events.zip
build/echo.zip:
	$(MAKE) -C cmd/echo OUT=$(root)/$@
build/square-events.zip:
	$(MAKE) -C cmd/square-events OUT=$(root)/$@

.PHONY: clean
clean: clean-echo clean-square-events
	rm -f template.yml
	rm -rf build/
.PHONY: clean-echo
clean-echo:
	$(MAKE) -C cmd/echo clean
.PHONY: clean-square-events
clean-square-events:
	$(MAKE) -C cmd/square-events clean
