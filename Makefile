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

test-update:
	go test ./pkg/cli/ -update-golden-files
	go test ./pkg/cobra/ -update-golden-files
	go test ./pkg/env/ -update-golden-files
	go test ./pkg/yaml/ -update-golden-files
