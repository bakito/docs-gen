# Include toolbox tasks
include ./.toolbox.mk

# Run go lint against code
lint: tb.golangci-lint
	$(TB_GOLANGCI_LINT) run --fix

# Run go mod tidy
tidy:
	go mod tidy

# Run ci tests
test:
	go test -cover -coverprofile coverage.out ./...
	go tool cover -func=coverage.out
