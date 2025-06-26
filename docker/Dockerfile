# Multi-stage Dockerfile for Go Get Images
# Supports cross-platform builds and testing

# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build for multiple platforms
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o go-get-imgs-linux-amd64 main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o go-get-imgs-linux-arm64 main.go
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o go-get-imgs-darwin-amd64 main.go
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -installsuffix cgo -o go-get-imgs-darwin-arm64 main.go
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -o go-get-imgs-windows-amd64.exe main.go
RUN CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -a -installsuffix cgo -o go-get-imgs-windows-arm64.exe main.go

# Test stage
FROM golang:1.23-alpine AS tester

# Install test dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy source code
COPY . .

# Run tests
RUN go test -v -race -cover
RUN go test -v -run TestIntegration
RUN go test -v -run TestCrossPlatform
RUN go test -bench=. -benchmem

# Runtime stage
FROM alpine:latest AS runtime

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/go-get-imgs-linux-amd64 .

# Make it executable
RUN chmod +x go-get-imgs-linux-amd64

# Create downloads directory
RUN mkdir -p downloads

# Expose volume for downloads
VOLUME ["/root/downloads"]

# Set entrypoint
ENTRYPOINT ["./go-get-imgs-linux-amd64"]
CMD ["--help"] 