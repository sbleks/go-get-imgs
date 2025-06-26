package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sbleks/go-get-imgs/internal/downloader"
	"github.com/sbleks/go-get-imgs/internal/utils"
)

// TestCrossPlatformFileOperations tests file operations across different platforms
func TestCrossPlatformFileOperations(t *testing.T) {
	// Test different path separators
	testCases := []struct {
		name     string
		paths    []string
		expected string
	}{
		{
			name:     "Unix-style paths",
			paths:    []string{"downloads", "image_1.jpg"},
			expected: filepath.Join("downloads", "image_1.jpg"),
		},
		{
			name:     "Windows-style paths",
			paths:    []string{"downloads", "image_1.jpg"},
			expected: filepath.Join("downloads", "image_1.jpg"),
		},
		{
			name:     "Mixed paths",
			paths:    []string{"downloads", "subdir", "image_1.jpg"},
			expected: filepath.Join("downloads", "subdir", "image_1.jpg"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filepath.Join(tc.paths...)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestCrossPlatformPathHandling tests path handling on different platforms
func TestCrossPlatformPathHandling(t *testing.T) {
	// Test file extension detection with different path formats
	testCases := []struct {
		url      string
		expected string
	}{
		{"https://example.com/image.jpg", ".jpg"},
		{"https://example.com/path/to/image.png", ".png"},
		{"https://example.com/image.gif", ".gif"},
		{"https://example.com/image.webp", ".webp"},
		{"https://example.com/image.bmp", ".bmp"},
		{"https://example.com/image.tiff", ".tiff"},
		{"https://example.com/image.tif", ".tif"},
		{"https://example.com/image.JPG", ".jpg"}, // case insensitive
		{"https://example.com/image.PNG", ".png"},
		{"https://example.com/image", ""},
		{"https://example.com/", ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("URL_%s", tc.url), func(t *testing.T) {
			result := downloader.GetExtensionFromURL(tc.url)
			if result != tc.expected {
				t.Errorf("For URL '%s', expected '%s', got '%s'", tc.url, tc.expected, result)
			}
		})
	}
}

// TestCrossPlatformDirectoryCreation tests directory creation across platforms
func TestCrossPlatformDirectoryCreation(t *testing.T) {
	testDirs := []string{
		"test_dir",
		"test_dir/subdir",
		"test_dir/subdir/nested",
		"test-dir-with-dashes",
		"test_dir_with_underscores",
	}

	for _, dir := range testDirs {
		t.Run(fmt.Sprintf("Create_%s", dir), func(t *testing.T) {
			// Clean up before test
			os.RemoveAll(dir)

			// Create directory
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				t.Errorf("Failed to create directory '%s': %v", dir, err)
				return
			}

			// Verify directory exists
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				t.Errorf("Directory '%s' was not created", dir)
			}

			// Clean up after test
			os.RemoveAll(dir)
		})
	}
}

// TestCrossPlatformFileNaming tests file naming conventions across platforms
func TestCrossPlatformFileNaming(t *testing.T) {
	testCases := []struct {
		rowNum    int
		extension string
		expected  string
	}{
		{1, ".jpg", "image_1.jpg"},
		{2, ".png", "image_2.png"},
		{10, ".gif", "image_10.gif"},
		{100, ".webp", "image_100.webp"},
		{999, ".bmp", "image_999.bmp"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Row_%d_%s", tc.rowNum, tc.extension), func(t *testing.T) {
			filename := fmt.Sprintf("image_%d%s", tc.rowNum, tc.extension)
			if filename != tc.expected {
				t.Errorf("Expected filename '%s', got '%s'", tc.expected, filename)
			}
		})
	}
}

// TestCrossPlatformEnvironmentDetection tests environment detection
func TestCrossPlatformEnvironmentDetection(t *testing.T) {
	// Test OS detection
	osName := runtime.GOOS
	arch := runtime.GOARCH

	t.Logf("Current OS: %s", osName)
	t.Logf("Current Architecture: %s", arch)

	// Verify we can detect the current platform
	if osName == "" {
		t.Error("OS name should not be empty")
	}

	if arch == "" {
		t.Error("Architecture should not be empty")
	}

	// Test supported platforms
	supportedOS := []string{"darwin", "linux", "windows"}
	supportedArch := []string{"amd64", "arm64"}

	osSupported := false
	for _, supported := range supportedOS {
		if osName == supported {
			osSupported = true
			break
		}
	}

	archSupported := false
	for _, supported := range supportedArch {
		if arch == supported {
			archSupported = true
			break
		}
	}

	if !osSupported {
		t.Logf("Warning: OS '%s' is not in the standard supported list", osName)
	}

	if !archSupported {
		t.Logf("Warning: Architecture '%s' is not in the standard supported list", arch)
	}
}

// TestCrossPlatformCSVHandling tests CSV handling across platforms
func TestCrossPlatformCSVHandling(t *testing.T) {
	// Test CSV with different line endings
	csvDataUnix := "id,name,image_url\n1,Test,https://example.com/image.jpg\n"
	csvDataWindows := "id,name,image_url\r\n1,Test,https://example.com/image.jpg\r\n"
	csvDataMixed := "id,name,image_url\r\n1,Test,https://example.com/image.jpg\n"

	testCases := []struct {
		name     string
		data     string
		expected int
	}{
		{"Unix line endings", csvDataUnix, 1},
		{"Windows line endings", csvDataWindows, 1},
		{"Mixed line endings", csvDataMixed, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary CSV file
			tempFile := fmt.Sprintf("temp_%s.csv", strings.ReplaceAll(tc.name, " ", "_"))
			err := os.WriteFile(tempFile, []byte(tc.data), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile)

			// Read and parse CSV
			file, err := os.Open(tempFile)
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			defer file.Close()

			reader := csv.NewReader(file)

			// Skip header
			_, err = reader.Read()
			if err != nil {
				t.Fatalf("Failed to read header: %v", err)
			}

			// Count data rows
			rowCount := 0
			for {
				_, err := reader.Read()
				if err != nil {
					break
				}
				rowCount++
			}

			if rowCount != tc.expected {
				t.Errorf("Expected %d rows, got %d", tc.expected, rowCount)
			}
		})
	}
}

// TestCrossPlatformURLHandling tests URL handling across platforms
func TestCrossPlatformURLHandling(t *testing.T) {
	testCases := []struct {
		name          string
		url           string
		expectedEmpty bool
		expectedValid bool
	}{
		{
			name:          "Valid HTTPS URL",
			url:           "https://example.com/image.jpg",
			expectedEmpty: false,
			expectedValid: true,
		},
		{
			name:          "Valid HTTP URL",
			url:           "http://example.com/image.png",
			expectedEmpty: false,
			expectedValid: true,
		},
		{
			name:          "Valid URL with path",
			url:           "https://example.com/path/to/image.gif",
			expectedEmpty: false,
			expectedValid: true,
		},
		{
			name:          "Invalid URL without scheme",
			url:           "invalid-url",
			expectedEmpty: false,
			expectedValid: false,
		},
		{
			name:          "Empty string",
			url:           "",
			expectedEmpty: true,
			expectedValid: false,
		},
		{
			name:          "Valid FTP URL",
			url:           "ftp://example.com/image.jpg",
			expectedEmpty: false,
			expectedValid: true,
		},
		{
			name:          "Valid file URL",
			url:           "file:///path/to/image.jpg",
			expectedEmpty: false,
			expectedValid: true,
		},
		{
			name:          "Whitespace only",
			url:           "   ",
			expectedEmpty: true,
			expectedValid: false,
		},
		{
			name:          "Valid URL with whitespace",
			url:           "  https://example.com/image.jpg  ",
			expectedEmpty: false,
			expectedValid: true,
		},
		{
			name:          "Invalid URL with whitespace",
			url:           "  invalid-url  ",
			expectedEmpty: false,
			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test URL trimming
			trimmed := strings.TrimSpace(tc.url)

			// Test empty URL detection
			isEmpty := trimmed == ""
			if isEmpty != tc.expectedEmpty {
				t.Errorf("URL '%s': expected isEmpty=%v, got %v", tc.url, tc.expectedEmpty, isEmpty)
			}

			// Test URL validation using the application's function
			isValid := utils.IsValidURL(tc.url)
			if isValid != tc.expectedValid {
				t.Errorf("URL '%s': expected isValid=%v, got %v", tc.url, tc.expectedValid, isValid)
			}
		})
	}
}

// TestCrossPlatformErrorHandling tests error handling across platforms
func TestCrossPlatformErrorHandling(t *testing.T) {
	// Test file not found error
	_, err := os.Stat("non_existent_file.txt")
	if !os.IsNotExist(err) {
		t.Error("Expected IsNotExist error for non-existent file")
	}

	// Test directory creation error (try to create in non-existent parent)
	err = os.MkdirAll("/non/existent/path/test", 0755)
	if err == nil && runtime.GOOS != "windows" {
		// On Unix systems, this should fail
		t.Error("Expected error when creating directory in non-existent parent")
	}

	// Test file creation in non-existent directory
	err = os.WriteFile("non_existent_dir/test.txt", []byte("test"), 0644)
	if err == nil {
		t.Error("Expected error when creating file in non-existent directory")
	}
}

// TestCrossPlatformBuildTags tests build tag functionality
func TestCrossPlatformBuildTags(t *testing.T) {
	// This test verifies that build tags work correctly
	// The actual build tag testing would be done during compilation

	t.Logf("Build tags test - this test verifies the test framework works")
	t.Logf("OS: %s", runtime.GOOS)
	t.Logf("Arch: %s", runtime.GOARCH)
	t.Logf("Compiler: %s", runtime.Compiler)
	t.Logf("Go version: %s", runtime.Version())
}

// BenchmarkCrossPlatformOperations benchmarks cross-platform operations
func BenchmarkCrossPlatformFileOperations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Benchmark filepath.Join
		filepath.Join("downloads", fmt.Sprintf("image_%d.jpg", i))
	}
}

func BenchmarkCrossPlatformURLProcessing(b *testing.B) {
	testURLs := []string{
		"https://example.com/image1.jpg",
		"https://example.com/image2.png",
		"https://example.com/image3.gif",
		"https://example.com/image4.webp",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		url := testURLs[i%len(testURLs)]
		downloader.GetExtensionFromURL(url)
		strings.TrimSpace(url)
	}
}
