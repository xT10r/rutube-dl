# Makefile for rutube-dl project

.PHONY: test test-unit test-integration test-ffmpeg test-utils test-i18n build clean bench coverage help

# Default target
all: test build

# Build the application
build:
	@echo "🔨 Building rutube-dl..."
	go build -o rutubedl.exe ./cmd/rutube-dl

# Run all tests
test: test-unit test-integration
	@echo "✅ All tests completed"

# Run unit tests
test-unit:
	@echo "🧪 Running unit tests..."
	go test -v ./pkg/...

# Run unit tests with short flag (skip network tests)
test-unit-short:
	@echo "🧪 Running unit tests (short mode)..."
	go test -short -v ./pkg/...

# Run integration tests
test-integration:
	@echo "🔧 Running integration tests..."
	go run ./test/integration.go

# Test specific packages
test-ffmpeg:
	@echo "🔧 Testing FFmpeg package..."
	go test -v ./pkg/ffmpeg

test-utils:
	@echo "🛠️  Testing Utils package..."
	go test -v ./pkg/utils

test-i18n:
	@echo "🌍 Testing i18n package..."
	go test -v ./pkg/i18n

# Run benchmarks
bench:
	@echo "📊 Running benchmarks..."
	go test -bench=. -benchmem ./pkg/...

# Generate test coverage
coverage:
	@echo "📈 Generating test coverage..."
	go test -coverprofile=coverage.out ./pkg/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Test FFmpeg download functionality
test-ffmpeg-download:
	@echo "⬇️  Testing FFmpeg download..."
	@echo "package main" > test_download.go
	@echo "import (\"fmt\"; \"github.com/StanislavKH/rutube-dl/pkg/ffmpeg\")" >> test_download.go
	@echo "func main() {" >> test_download.go
	@echo "  mgr := ffmpeg.NewFFmpegManager()" >> test_download.go
	@echo "  err := mgr.EnsureFFmpeg(false)" >> test_download.go
	@echo "  if err != nil { fmt.Printf(\"Error: %v\\n\", err) } else { fmt.Println(\"✅ FFmpeg ready!\") }" >> test_download.go
	@echo "}" >> test_download.go
	go run test_download.go
	rm test_download.go

# Clean build artifacts and test files
clean:
	@echo "🧹 Cleaning..."
	rm -f rutubedl.exe
	rm -f coverage.out coverage.html
	rm -f test_download.go

# Run tests with race detection
test-race:
	@echo "🏃 Running tests with race detection..."
	go test -race -v ./pkg/...

# Lint the code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Format the code
fmt:
	@echo "💫 Formatting code..."
	go fmt ./...

# Tidy dependencies
tidy:
	@echo "📦 Tidying dependencies..."
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  test               - Run all tests"
	@echo "  test-unit          - Run unit tests"
	@echo "  test-unit-short    - Run unit tests (skip network tests)"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-ffmpeg        - Test FFmpeg package"
	@echo "  test-utils         - Test utils package"
	@echo "  test-i18n          - Test i18n package"
	@echo "  test-ffmpeg-download - Test FFmpeg download functionality"
	@echo "  bench              - Run benchmarks"
	@echo "  coverage           - Generate test coverage report"
	@echo "  test-race          - Run tests with race detection"
	@echo "  lint               - Lint the code"
	@echo "  fmt                - Format the code"
	@echo "  tidy               - Tidy dependencies"
	@echo "  clean              - Clean build artifacts"
	@echo "  help               - Show this help"