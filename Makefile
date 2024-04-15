APP_NAME = digitalvenue
ENVIRONMENT = dev
CODE_BUCKET = $(APP_NAME)-codebucket
ROOT = $(shell git rev-parse --show-toplevel)
SERVICES = $(shell ls functions)

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
	if [ -z "$(ROUTE_53_HOSTED_ZONE_ID)" ]; then \
		echo "ROUTE_53_HOSTED_ZONE_ID is not set"; \
		exit 1; \
	fi
	aws cloudformation deploy \
		--stack-name $(APP_NAME)-$(ENVIRONMENT) \
		--template-file $(ROOT)/template.yml \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			Route53HostedZoneId=${ROUTE_53_HOSTED_ZONE_ID} \
			Environment=$(ENVIRONMENT) \
			CodeBucketName=$(CODE_BUCKET)

.PHONY: build
build: $(addprefix build/, $(addsuffix .zip, $(SERVICES)))
build/%.zip: functions/%
	$(MAKE) -C $< OUT=$(ROOT)/$@

.PHONY: test
test: build
	go test ./...
	go vet ./...

.PHONY: clean
clean: $(addprefix clean-, $(SERVICES))
	rm -f template.yml
	rm -rf build/
clean-%:
	$(MAKE) -C functions/$* clean
