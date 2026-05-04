.PHONY: run gen-docs mock

run:
	go run ./cmd/server

gen-docs:
	go run ./cmd/gen-docs

mock-gen:
	mockery

.PHONY: tidy fmt lint vulncheck test test-cover test-clean pre-commit

tidy:
	go mod tidy

fmt:
	golangci-lint fmt ./...

lint:
	golangci-lint run ./...

vulncheck:
	govulncheck ./...

test:
	go test ./... -v

test-cover:
	go test ./... -coverprofile=coverage.out
	grep -vE "mocks/|\.gen\.go" coverage.out > coverage.nomocks.out
	go tool cover -func=coverage.nomocks.out

test-clean:
	go clean -testcache && go test ./... -v -cover

pre-commit: tidy fmt lint gen-docs test-cover vulncheck
	@echo =============================
	@echo All pre-commit checks passed.
