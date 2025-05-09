.PHONY: build
build:
	sqlc generate
	go build -o build/dvserver ./cmd/dvserver

.PHONY: run-server
run-server: build
	go run ./cmd/server

.PHONY: test
test: build
	go test ./...
	golangci-lint run

.PHONY: migrate
migrate:
	goose postgres "$(DV_DATABASE_URL)" up
