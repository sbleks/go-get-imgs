# Go Get Images

A cross-platform Go application that reads a CSV file row by row and downloads images from URLs in a specified column. Supports Windows, macOS, and Linux.

## Features

- **Cross-Platform Support**: Runs on Windows, macOS, and Linux (amd64 and arm64)
- Reads CSV files row by row
- Downloads images from URLs in any specified column
- Automatically detects image file types from content-type headers
- Creates a `downloads` directory to store images
- Provides detailed progress and error reporting
- Handles various image formats (JPG, PNG, GIF, WebP, BMP, TIFF)
- Includes timeout protection for HTTP requests
- Comprehensive test suite with unit and integration tests
- Cross-platform CI/CD with GitHub Actions
- Docker support for containerized builds and testing

## Supported Platforms

| OS | Architecture | Status |
|---|---|---|
| Windows | amd64, arm64 | ✅ Supported |
| macOS | amd64, arm64 | ✅ Supported |
| Linux | amd64, arm64 | ✅ Supported |

## Requirements

- Go 1.23.0 or later
- Windows, macOS, or Linux

## Quick Start

### Using Pre-built Binaries

Download the appropriate binary for your platform from the [releases page](https://github.com/your-repo/go-get-imgs/releases).

### Building from Source

```bash
# Clone the repository
git clone https://github.com/your-repo/go-get-imgs.git
cd go-get-imgs

# Install dependencies
go mod tidy

# Build for current platform
go build -o go-get-imgs main.go

# Run the application
./go-get-imgs sample.csv 3
```

## Cross-Platform Builds

### Using Makefile (macOS/Linux)

```bash
# Build for all platforms
make build-all

# Build for specific platform
make build-darwin    # macOS
make build-linux     # Linux
make build-windows   # Windows

# Create distribution packages
make dist

# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Using Windows Batch File

```cmd
# Run the interactive build menu
build-all.bat

# Or use specific commands
go build -o go-get-imgs.exe main.go
```

### Using Docker

```bash
# Build and test in Docker
docker-compose up builder
docker-compose up tester

# Run the application in Docker
docker-compose up runtime

# Development environment
docker-compose run dev
```

### Manual Cross-Platform Builds

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o go-get-imgs-darwin-amd64 main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o go-get-imgs-darwin-arm64 main.go

# Linux (Intel)
GOOS=linux GOARCH=amd64 go build -o go-get-imgs-linux-amd64 main.go

# Linux (ARM)
GOOS=linux GOARCH=arm64 go build -o go-get-imgs-linux-arm64 main.go

# Windows (Intel)
GOOS=windows GOARCH=amd64 go build -o go-get-imgs-windows-amd64.exe main.go

# Windows (ARM)
GOOS=windows GOARCH=arm64 go build -o go-get-imgs-windows-arm64.exe main.go
```

## Usage

### Basic Usage

```bash
go-get-imgs <csv-file> <url-column-index>
```

### Examples

```bash
# Download images from column 3
./go-get-imgs sample.csv 3

# Download images from column 2
./go-get-imgs data.csv 2

# On Windows
go-get-imgs.exe sample.csv 3
```

## CSV Format

Your CSV file should have at least the number of columns specified by `url-column-index`, with image URLs in that column:

```csv
id,name,image_url,description
1,Image 1,https://example.com/image1.jpg,First image
2,Image 2,https://example.com/image2.png,Second image
```

## Output

- Images are downloaded to a `downloads` directory
- Files are named as `image_1.jpg`, `image_2.png`, etc. (based on row number)
- The application shows progress and provides a summary at the end

## Error Handling

The application handles various error scenarios:
- Missing or invalid CSV files
- Network timeouts (30-second timeout)
- Invalid URLs
- HTTP errors
- File system errors
- Empty URLs in specified column

## Testing

This project includes a comprehensive test suite with unit tests, integration tests, and cross-platform tests.

### Running Tests

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run specific test categories
go test -v -run TestDownloadImage
go test -v -run TestIntegration
go test -v -run TestCrossPlatform

# Run benchmarks
go test -bench=. -benchmem

# Using test scripts
./run-tests.sh          # macOS/Linux
run-tests.bat           # Windows
```

### Test Coverage

The test suite covers:
- ✅ CSV file parsing and validation
- ✅ Image downloading with different content types
- ✅ Error handling (network errors, timeouts, invalid URLs)
- ✅ File extension detection from content-type headers
- ✅ File extension detection from URLs
- ✅ Command line argument validation
- ✅ File and directory operations
- ✅ Complete application workflows
- ✅ Large file processing
- ✅ Error scenarios and edge cases
- ✅ Cross-platform compatibility
- ✅ Path handling on different operating systems

### Test Structure

- **`main_test.go`** - Unit tests for individual functions
- **`integration_test.go`** - Integration tests for complete workflows
- **`cross_platform_test.go`** - Cross-platform compatibility tests
- **`test_helpers.go`** - Test utilities and helper functions

## Continuous Integration

The project includes GitHub Actions workflows for:

- **Cross-platform testing** on Windows, macOS, and Linux
- **Code quality checks** (linting, formatting, security)
- **Automated releases** with pre-built binaries
- **Docker builds** and testing

### CI/CD Pipeline

1. **Test Matrix**: Runs tests on all supported platforms
2. **Code Quality**: Linting, formatting, and security checks
3. **Integration Tests**: Complete workflow testing
4. **Cross-platform Tests**: Platform-specific functionality
5. **Release Builds**: Automated binary creation for releases

## Docker Support

### Building with Docker

```bash
# Build all platforms
docker build --target builder .

# Run tests
docker build --target tester .

# Run the application
docker build --target runtime .
docker run -v $(pwd)/downloads:/root/downloads -v $(pwd)/sample.csv:/root/sample.csv go-get-imgs
```

### Using Docker Compose

```bash
# Build all platforms and create distribution packages
docker-compose up builder

# Run all tests
docker-compose up tester

# Run the application
docker-compose up runtime

# Development environment
docker-compose run dev
```

## Development

### Project Structure

```
go-get-imgs/
├── main.go                 # Main application
├── main_test.go           # Unit tests
├── integration_test.go    # Integration tests
├── cross_platform_test.go # Cross-platform tests
├── test_helpers.go        # Test utilities
├── Makefile              # Build system (macOS/Linux)
├── build-all.bat         # Build system (Windows)
├── Dockerfile            # Docker configuration
├── docker-compose.yml    # Docker Compose configuration
├── .github/workflows/    # CI/CD workflows
├── sample.csv            # Sample data
└── README.md             # This file
```

### Development Workflow

1. **Setup**: `go mod tidy`
2. **Test**: `go test -v`
3. **Build**: `make build` or `go build`
4. **Cross-platform**: `make build-all`
5. **Docker**: `docker-compose up dev`

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Security check (requires gosec)
gosec ./...
```

## Sample Data

A sample CSV file (`sample.csv`) is included with test image URLs from Picsum Photos.

## Dependencies

- Standard Go libraries (no external dependencies required)
- Uses `encoding/csv` for CSV parsing
- Uses `net/http` for image downloads
- Uses `os` and `path/filepath` for file operations

## License

This project is open source and available under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## Support

For issues and questions:
- Check the [GitHub Issues](https://github.com/your-repo/go-get-imgs/issues)
- Review the test suite for usage examples
- Check the sample CSV file for format examples 