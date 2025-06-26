@echo off
setlocal enabledelayedexpansion

REM Go Get Images - Cross-Platform Build System for Windows
REM Supports macOS, Linux, and Windows

REM Variables
set BINARY_NAME=go-get-imgs
set BUILD_DIR=build
set DIST_DIR=dist

REM Get version from git or use default
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=dev

REM Get build time
for /f "tokens=1-3 delims=/ " %%a in ('date /t') do set BUILD_DATE=%%a-%%b-%%c
for /f "tokens=1-2 delims=: " %%a in ('time /t') do set BUILD_TIME=%%a-%%b

REM Get git commit
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

REM Build flags
set LDFLAGS=-ldflags "-X main.Version=%VERSION% -X main.BuildTime=%BUILD_DATE%_%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%"

echo === Go Get Images Cross-Platform Build System ===
echo Version: %VERSION%
echo Build Time: %BUILD_DATE%_%BUILD_TIME%
echo Git Commit: %GIT_COMMIT%
echo.

:menu
echo Select an option:
echo 1. Build for current platform (Windows)
echo 2. Build for all platforms
echo 3. Build for macOS only
echo 4. Build for Linux only
echo 5. Build for Windows only
echo 6. Run tests
echo 7. Run tests with coverage
echo 8. Run benchmarks
echo 9. Create distribution packages
echo 10. Clean build artifacts
echo 11. Install dependencies
echo 12. Format code
echo 13. Vet code
echo 14. Run all checks
echo 15. Exit
echo.
set /p choice="Enter your choice (1-15): "

if "%choice%"=="1" goto build-current
if "%choice%"=="2" goto build-all-platforms
if "%choice%"=="3" goto build-darwin
if "%choice%"=="4" goto build-linux
if "%choice%"=="5" goto build-windows
if "%choice%"=="6" goto run-tests
if "%choice%"=="7" goto run-tests-coverage
if "%choice%"=="8" goto run-benchmarks
if "%choice%"=="9" goto create-dist
if "%choice%"=="10" goto clean
if "%choice%"=="11" goto deps
if "%choice%"=="12" goto fmt
if "%choice%"=="13" goto vet
if "%choice%"=="14" goto check
if "%choice%"=="15" goto exit
echo Invalid choice. Please try again.
echo.
goto menu

:build-current
echo Building for current platform (Windows)...
if not exist %BUILD_DIR% mkdir %BUILD_DIR%
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%.exe main.go
if %ERRORLEVEL% EQU 0 (
    echo Build successful! Created %BUILD_DIR%\%BINARY_NAME%.exe
) else (
    echo Build failed!
)
goto end

:build-all-platforms
echo Building for all platforms...
if not exist %BUILD_DIR% mkdir %BUILD_DIR%

echo Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-darwin-amd64 main.go

echo Building for macOS (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-darwin-arm64 main.go

echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-linux-amd64 main.go

echo Building for Linux (arm64)...
set GOOS=linux
set GOARCH=arm64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-linux-arm64 main.go

echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-windows-amd64.exe main.go

echo Building for Windows (arm64)...
set GOOS=windows
set GOARCH=arm64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-windows-arm64.exe main.go

echo All builds completed!
goto end

:build-darwin
echo Building for macOS...
if not exist %BUILD_DIR% mkdir %BUILD_DIR%
set GOOS=darwin
set GOARCH=amd64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-darwin-amd64 main.go
set GOOS=darwin
set GOARCH=arm64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-darwin-arm64 main.go
echo macOS builds completed!
goto end

:build-linux
echo Building for Linux...
if not exist %BUILD_DIR% mkdir %BUILD_DIR%
set GOOS=linux
set GOARCH=amd64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-linux-amd64 main.go
set GOOS=linux
set GOARCH=arm64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-linux-arm64 main.go
echo Linux builds completed!
goto end

:build-windows
echo Building for Windows...
if not exist %BUILD_DIR% mkdir %BUILD_DIR%
set GOOS=windows
set GOARCH=amd64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-windows-amd64.exe main.go
set GOOS=windows
set GOARCH=arm64
go build %LDFLAGS% -o %BUILD_DIR%\%BINARY_NAME%-windows-arm64.exe main.go
echo Windows builds completed!
goto end

:run-tests
echo Running tests...
go test -v -race -cover
goto end

:run-tests-coverage
echo Running tests with coverage...
go test -v -race -coverprofile=coverage.out
if exist coverage.out (
    go tool cover -func=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    echo Coverage report generated: coverage.html
    del coverage.out
)
goto end

:run-benchmarks
echo Running benchmarks...
go test -bench=. -benchmem
goto end

:create-dist
echo Creating distribution packages...
if not exist %DIST_DIR% mkdir %DIST_DIR%

echo Creating Windows packages...
if exist %BUILD_DIR%\%BINARY_NAME%-windows-amd64.exe (
    powershell -Command "Compress-Archive -Path '%BUILD_DIR%\%BINARY_NAME%-windows-amd64.exe' -DestinationPath '%DIST_DIR%\%BINARY_NAME%-windows-amd64-%VERSION%.zip'"
)
if exist %BUILD_DIR%\%BINARY_NAME%-windows-arm64.exe (
    powershell -Command "Compress-Archive -Path '%BUILD_DIR%\%BINARY_NAME%-windows-arm64.exe' -DestinationPath '%DIST_DIR%\%BINARY_NAME%-windows-arm64-%VERSION%.zip'"
)

echo Distribution packages created in %DIST_DIR% directory!
goto end

:clean
echo Cleaning build artifacts...
if exist %BUILD_DIR% rmdir /s /q %BUILD_DIR%
if exist %DIST_DIR% rmdir /s /q %DIST_DIR%
go clean -cache
echo Cleanup completed!
goto end

:deps
echo Installing dependencies...
go mod tidy
go mod download
echo Dependencies installed!
goto end

:fmt
echo Formatting code...
go fmt ./...
echo Code formatting completed!
goto end

:vet
echo Vetting code...
go vet ./...
echo Code vetting completed!
goto end

:check
echo Running all checks...
echo Formatting code...
go fmt ./...
echo Vetting code...
go vet ./...
echo Running tests...
go test -v -race -cover
echo All checks completed!
goto end

:end
echo.
echo Press any key to return to menu...
pause >nul
goto menu

:exit
echo Goodbye!
exit /b 0 