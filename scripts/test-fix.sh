#!/bin/bash

echo "Testing URL validation fixes..."

# Test the isValidURL function specifically
echo "Running isValidURL tests..."
go test -v -run TestIsValidURL

echo ""
echo "Running cross-platform URL handling tests..."
go test -v -run TestCrossPlatformURLHandling

echo ""
echo "Running integration tests..."
go test -v -run TestIntegrationWithErrors

echo ""
echo "All tests completed!" 