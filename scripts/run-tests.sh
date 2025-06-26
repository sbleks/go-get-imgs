#!/bin/bash

echo "=== Go Get Images Test Suite ==="
echo

# Function to run tests with different options
run_tests() {
    echo "Running: $1"
    echo "Command: $2"
    echo "---"
    eval "$2"
    echo "---"
    echo
}

# Run all tests
run_tests "All Tests" "go test -v"

# Run tests with coverage
run_tests "Tests with Coverage" "go test -cover"

# Run specific test categories
run_tests "Unit Tests Only" "go test -v -run 'Test(DownloadImage|GetExtension|CSVProcessing|CommandLine|FileOperations)'"

run_tests "Integration Tests Only" "go test -v -run 'TestIntegration'"

# Run benchmarks
run_tests "Benchmarks" "go test -bench=."

# Generate coverage report
echo "Generating detailed coverage report..."
go test -coverprofile=coverage.out
if [ -f coverage.out ]; then
    echo "Coverage by function:"
    go tool cover -func=coverage.out
    echo
    echo "HTML coverage report generated: coverage.html"
    go tool cover -html=coverage.out -o coverage.html
    rm coverage.out
fi

echo "=== Test Suite Complete ===" 