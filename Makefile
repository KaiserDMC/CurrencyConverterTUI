APP_NAME=currency-converter
BUILD_DIR=bin
CMD_DIR=cmd/CurrencyConverterTUI

# Set the Go environment variables
GO=go
GO_BUILD=$(GO) build
GO_RUN=$(GO) run
GO_MOD_TIDY=$(GO) mod tidy

# Build the application
build:
	@echo "Building the application..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	$(GO_BUILD) -o $(BUILD_DIR)/$(APP_NAME).exe ./$(CMD_DIR)/main.go
	$(GO_BUILD) -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)/main.go

# Run the application
run:
	$(GO_RUN) ./$(CMD_DIR)/main.go

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)

# Install dependencies
deps:
	$(GO_MOD_TIDY)

# Format code
fmt:
	$(GO) fmt ./...

# Lint code (optional, if you have `golangci-lint` installed)
lint:
	golangci-lint run

# Build and run the application
all: build run
