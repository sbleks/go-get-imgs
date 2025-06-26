#!/bin/bash

# Go Get Images - Run Script
# Simple script to build and run the application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Go Get Images - Build and Run${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    exit 1
fi

# Build the application
echo -e "${YELLOW}Building application...${NC}"
go build -o go-get-imgs cmd/go-get-imgs/main.go

# Check if build was successful
if [ ! -f "go-get-imgs" ]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}Build successful!${NC}"

# Check if CSV file is provided
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}No arguments provided. Using sample data...${NC}"
    echo -e "${GREEN}Usage: $0 <csv-file> <url-column-index>${NC}"
    echo -e "${GREEN}Example: $0 examples/sample.csv 3${NC}"
    echo ""
    echo -e "${YELLOW}Running with sample data...${NC}"
    ./go-get-imgs examples/sample.csv 3
else
    echo -e "${GREEN}Running with provided arguments...${NC}"
    ./go-get-imgs "$@"
fi

echo -e "${GREEN}Done!${NC}" 