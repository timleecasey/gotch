# Makefile for gotch - PyTorch 2.10.0 Go bindings
# macOS with MPS support

# Path to libtorch installation
LIBTORCH_PATH := /Users/tlc/src/gotch/libtorch-2.10.0-macos

# Environment variables for building and testing
export CPATH := $(LIBTORCH_PATH)/include/torch/csrc/api/include:$(LIBTORCH_PATH)/include
export LIBRARY_PATH := $(LIBTORCH_PATH)/lib
export DYLD_LIBRARY_PATH := $(LIBTORCH_PATH)/lib

# Go test flags
TEST_FLAGS := -v
TEST_TIMEOUT := 5m

.PHONY: all build test test-nn test-ts test-all clean help

# Default target
all: build

# Build all core packages
build:
	@echo "Building gotch core packages..."
	@go build -v . ./ts ./nn ./vision

# Run all tests
test-all: test-nn test-ts
	@echo "All tests completed"

# Run nn package tests
test-nn:
	@echo "Running nn package tests..."
	@go test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./nn

# Run ts package tests
test-ts:
	@echo "Running ts package tests..."
	@go test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./ts

# Run specific test in nn package
# Usage: make test-nn-specific TEST=TestInitTensor_Memcheck
test-nn-specific:
	@echo "Running specific nn test: $(TEST)..."
	@go test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) -run $(TEST) ./nn

# Run specific test in ts package
# Usage: make test-ts-specific TEST=TestTensor
test-ts-specific:
	@echo "Running specific ts test: $(TEST)..."
	@go test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) -run $(TEST) ./ts

# Run tests with coverage
test-nn-coverage:
	@echo "Running nn tests with coverage..."
	@go test -v -timeout $(TEST_TIMEOUT) -coverprofile=coverage-nn.out ./nn
	@go tool cover -html=coverage-nn.out -o coverage-nn.html
	@echo "Coverage report saved to coverage-nn.html"

test-ts-coverage:
	@echo "Running ts tests with coverage..."
	@go test -v -timeout $(TEST_TIMEOUT) -coverprofile=coverage-ts.out ./ts
	@go tool cover -html=coverage-ts.out -o coverage-ts.html
	@echo "Coverage report saved to coverage-ts.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@go clean -cache -testcache
	@rm -f coverage-*.out coverage-*.html

# Display build environment
env:
	@echo "Build Environment:"
	@echo "  LIBTORCH_PATH:      $(LIBTORCH_PATH)"
	@echo "  CPATH:              $(CPATH)"
	@echo "  LIBRARY_PATH:       $(LIBRARY_PATH)"
	@echo "  DYLD_LIBRARY_PATH:  $(DYLD_LIBRARY_PATH)"
	@echo ""
	@echo "Go version:"
	@go version
	@echo ""
	@echo "LibTorch version: 2.10.0"
	@echo "MPS support: Available on Apple Silicon"

# Check if MPS is available
check-mps:
	@echo "Checking MPS availability..."
	@go run -exec 'env DYLD_LIBRARY_PATH=$(DYLD_LIBRARY_PATH)' tools/check_device.go || echo "Create tools/check_device.go to test device availability"

# Help target
help:
	@echo "Gotch Makefile - PyTorch 2.10.0 Go Bindings"
	@echo ""
	@echo "Targets:"
	@echo "  make build              - Build all core packages"
	@echo "  make test-all           - Run all tests (nn + ts)"
	@echo "  make test-nn            - Run nn package tests"
	@echo "  make test-ts            - Run ts package tests"
	@echo "  make test-nn-specific   - Run specific nn test (TEST=TestName)"
	@echo "  make test-ts-specific   - Run specific ts test (TEST=TestName)"
	@echo "  make test-nn-coverage   - Run nn tests with coverage report"
	@echo "  make test-ts-coverage   - Run ts tests with coverage report"
	@echo "  make clean              - Clean build artifacts and caches"
	@echo "  make env                - Display build environment"
	@echo "  make check-mps          - Check MPS device availability"
	@echo "  make help               - Show this help message"
	@echo ""
	@echo "Environment:"
	@echo "  LibTorch: $(LIBTORCH_PATH)"
	@echo "  PyTorch version: 2.10.0"
	@echo "  Platform: macOS with MPS support"
	@echo ""
	@echo "Examples:"
	@echo "  make test-nn"
	@echo "  make test-nn-specific TEST=TestInitTensor_Memcheck"
	@echo "  make test-all"
