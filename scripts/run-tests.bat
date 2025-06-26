@echo off
echo === Go Get Images Test Suite ===
echo.

REM Run all tests
echo Running: All Tests
echo Command: go test -v
echo ---
go test -v
echo ---
echo.

REM Run tests with coverage
echo Running: Tests with Coverage
echo Command: go test -cover
echo ---
go test -cover
echo ---
echo.

REM Run unit tests only
echo Running: Unit Tests Only
echo Command: go test -v -run "Test(DownloadImage^|GetExtension^|CSVProcessing^|CommandLine^|FileOperations)"
echo ---
go test -v -run "Test(DownloadImage|GetExtension|CSVProcessing|CommandLine|FileOperations)"
echo ---
echo.

REM Run integration tests only
echo Running: Integration Tests Only
echo Command: go test -v -run "TestIntegration"
echo ---
go test -v -run "TestIntegration"
echo ---
echo.

REM Run benchmarks
echo Running: Benchmarks
echo Command: go test -bench=.
echo ---
go test -bench=.
echo ---
echo.

REM Generate coverage report
echo Generating detailed coverage report...
go test -coverprofile=coverage.out
if exist coverage.out (
    echo Coverage by function:
    go tool cover -func=coverage.out
    echo.
    echo HTML coverage report generated: coverage.html
    go tool cover -html=coverage.out -o coverage.html
    del coverage.out
)

echo === Test Suite Complete ===
pause 