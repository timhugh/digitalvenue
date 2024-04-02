APP_NAME = digitalvenue
ENVIRONMENT = dev
CODE_BUCKET = $(APP_NAME)-codebucket
ROOT = $(shell git rev-parse --show-toplevel)

.PHONY: validate
validate:
	aws cloudformation validate-template --template-body file://$(ROOT)/.cloudformation/env.yml

.PHONY: codebucket
codebucket:
	aws cloudformation deploy --stack-name $(CODE_BUCKET) \
		--template-file $(ROOT)/.cloudformation/bucket.yml \
		--parameter-overrides \
			BucketName=$(CODE_BUCKET)

.PHONY: package
package: codebucket build
	aws cloudformation package \
		--template-file $(ROOT)/.cloudformation/env.yml \
		--output-template-file $(ROOT)/template.yml \
		--s3-bucket $(CODE_BUCKET) \
		--s3-prefix $(APP_NAME)/$(ENVIRONMENT)

.PHONY: deploy
deploy: package
	aws cloudformation deploy \
		--stack-name $(APP_NAME)-$(ENVIRONMENT) \
		--template-file $(ROOT)/template.yml \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			Route53HostedZoneId=${ROUTE_53_HOSTED_ZONE_ID} \
			Environment=$(ENVIRONMENT) \
			CodeBucketName=$(CODE_BUCKET)

.PHONY: build
build: build/echo.zip build/square-events.zip
build/echo.zip:
	$(MAKE) -C cmd/echo OUT=$(ROOT)/$@
build/square-events.zip:
	$(MAKE) -C cmd/square-events OUT=$(ROOT)/$@

.PHONY: test
test: build
	go test ./...
	go vet ./...

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
