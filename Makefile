APP_NAME = digitalvenue
ENVIRONMENT = dev
CODE_BUCKET = $(APP_NAME)-codebucket
ROOT = $(shell git rev-parse --show-toplevel)
SERVICES = $(shell ls functions)

default: build

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
	if [ -z "$(PAPERTRAIL_LOG_PUSH_URL)" ]; then \
		echo "PAPERTRAIL_LOG_PUSH_URL is not set"; \
		exit 1; \
	fi
	aws cloudformation deploy \
		--stack-name $(APP_NAME)-$(ENVIRONMENT) \
		--template-file $(ROOT)/template.yml \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			Route53HostedZoneId=${ROUTE_53_HOSTED_ZONE_ID} \
			Environment=$(ENVIRONMENT) \
			CodeBucketName=$(CODE_BUCKET) \
			PapertrailLogPushURL=$(PAPERTRAIL_LOG_PUSH_URL)

.PHONY: build
build: $(addprefix build/, $(addsuffix .zip, $(SERVICES)))
build/%.zip: build/%/bootstrap
	zip -rj $@ $^

build/%/bootstrap: functions/% functions/%/wire_gen.go util
	GOOS=linux GOARCH=arm64 go build -o $@ ./functions/$*

.PHONY: wire
wire: $(addprefix functions/, $(addsuffix /wire_gen.go, $(SERVICES)))
functions/%/wire_gen.go:
	if [ -f ./functions/$*/wire.go ]; then wire gen ./functions/$*; fi

.PHONY: test
test: build
	go test ./...

.PHONY: lint
lint: test
	golangci-lint run

.PHONY: clean
clean: $(addprefix clean-, $(SERVICES))
	rm -f template.yml
	rm -rf build/
clean-%:
	rm -f functions/$*/wire_gen.go
