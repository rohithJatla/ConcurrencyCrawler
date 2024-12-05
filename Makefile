# Default target
all: build

# Define variables
BINARY_NAME := crawler
PACKAGE_DIR := ./pkg/crawler

# Define PHONY targets
.PHONY: all build test lint run clean install-lint

# Build the project
build:
	@echo 'Building the project...'
	@go build -o bin/$(BINARY_NAME) $(PACKAGE_DIR)

# Run tests
test:
	@echo "Running tests..."
	@go test -v $(PACKAGE_DIR)/...

# Lint the code
lint:tool/.golangci-lint.$(GOLANGCI_LINT_VERSION)
	@echo "Linting the code..."@tool/golangci-lint run $(PACKAGE_DIR)/...


# Fix target (auto-fix lint issues)
fix:tool/.golangci-lint.$(GOLANGCI_LINT_VERSION)
	@echo "Fixing lint issues..."	@tool/golangci-lint run --fix $(PACKAGE_DIR)/...

# Run the program
run:
	build	@echo "Running the program..."  |  @./bin/$(BINARY_NAME)

# Clean up generated files
clean:
	@echo "Cleaning up..." | @rm -rf bin/$(BINARY_NAME)

# Install golangci-lint locally if not available
GOLANGCI_LINT := $(shell which golangci-lint)

ifeq ($(GOLANGCI_LINT),)
lint: install-lint
endif

install-lint:
	@echo "Installing golangci-lint..." | @curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s --

# Ensure tool is built
tool: tool/golangci-lint

tool/golangci-lint: tool/.golangci-lint.$(GOLANGCI_LINT_VERSION)
	@mkdir -p tool ||@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@GOBIN="$(PWD)/tool" go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

tool/.golangci-lint.$(GOLANGCI_LINT_VERSION):
	@rm -f tool/.golangci-lint.*
	@mkdir -p tool
	@touch $@
