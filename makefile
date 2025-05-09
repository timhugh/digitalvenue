.PHONY: build
build:
	sqlc generate
	go build -o build/ ./...

.PHONY: run-server
run-server: build
	go run ./cmd/dvserver

.PHONY: test
test: build
	go test ./...
	golangci-lint run

.PHONY: migrate
migrate:
	goose postgres "$(DV_DATABASE_URL)" up

.PHONY: clean
clean:
	rm -rf build/
	rm -rf db/
