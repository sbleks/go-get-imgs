package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sbleks/go-get-imgs/internal/downloader"
	"github.com/sbleks/go-get-imgs/internal/utils"
)

// TestMain sets up and tears down test environment
func TestMain(m *testing.M) {
	// Create test downloads directory
	os.MkdirAll("test_downloads", 0755)

	// Run tests
	code := m.Run()

	// Cleanup
	os.RemoveAll("test_downloads")

	os.Exit(code)
}

// TestDownloadImage tests the downloadImage function
func TestDownloadImage(t *testing.T) {
	// Create a test server that serves a simple image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte("fake image data"))
	}))
	defer server.Close()

	d := downloader.NewDownloader(30 * time.Second)

	// Test successful download
	err := d.DownloadImage(server.URL, "test_downloads", 1)
	if err != nil {
		t.Errorf("Expected successful download, got error: %v", err)
	}

	// Check if file was created
	expectedFile := filepath.Join("test_downloads", "image_1.jpg")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to be created", expectedFile)
	}

	// Test with different content types
	testCases := []struct {
		contentType string
		expectedExt string
	}{
		{"image/png", ".png"},
		{"image/gif", ".gif"},
		{"image/webp", ".webp"},
		{"image/bmp", ".bmp"},
		{"image/tiff", ".tiff"},
		{"unknown/type", ".jpg"}, // fallback
	}

	for i, tc := range testCases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", tc.contentType)
			w.Write([]byte("fake image data"))
		}))

		err := d.DownloadImage(server.URL, "test_downloads", i+100)
		if err != nil {
			t.Errorf("Test case %d: Expected successful download, got error: %v", i, err)
		}

		expectedFile := filepath.Join("test_downloads", fmt.Sprintf("image_%d%s", i+100, tc.expectedExt))
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Test case %d: Expected file %s to be created", i, expectedFile)
		}

		server.Close()
	}
}

// TestDownloadImageErrors tests error scenarios
func TestDownloadImageErrors(t *testing.T) {
	d := downloader.NewDownloader(30 * time.Second)
	// Test invalid URL
	err := d.DownloadImage("invalid-url", "test_downloads", 1)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}

	// Test server returning error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	err = d.DownloadImage(server.URL, "test_downloads", 1)
	if err == nil {
		t.Error("Expected error for 404 status, got nil")
	}

	// Test server timeout
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(35 * time.Second)
		w.Write([]byte("too late"))
	}))
	defer slowServer.Close()

	err = d.DownloadImage(slowServer.URL, "test_downloads", 1)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

// TestGetExtensionFromContentType tests content type detection
func TestGetExtensionFromContentType(t *testing.T) {
	testCases := []struct {
		contentType string
		expected    string
	}{
		{"image/jpeg", ".jpg"},
		{"image/jpg", ".jpg"},
		{"image/png", ".png"},
		{"image/gif", ".gif"},
		{"image/webp", ".webp"},
		{"image/bmp", ".bmp"},
		{"image/tiff", ".tiff"},
		{"text/html", ""},
		{"application/json", ""},
		{"", ""},
	}

	for _, tc := range testCases {
		result := downloader.GetExtensionFromContentType(tc.contentType)
		if result != tc.expected {
			t.Errorf("For content type '%s', expected '%s', got '%s'", tc.contentType, tc.expected, result)
		}
	}
}

// TestGetExtensionFromURL tests URL extension detection
func TestGetExtensionFromURL(t *testing.T) {
	testCases := []struct {
		url      string
		expected string
	}{
		{"https://example.com/image.jpg", ".jpg"},
		{"https://example.com/image.jpeg", ".jpeg"},
		{"https://example.com/image.png", ".png"},
		{"https://example.com/image.gif", ".gif"},
		{"https://example.com/image.webp", ".webp"},
		{"https://example.com/image.bmp", ".bmp"},
		{"https://example.com/image.tiff", ".tiff"},
		{"https://example.com/image.tif", ".tif"},
		{"https://example.com/image.JPG", ".jpg"}, // case insensitive
		{"https://example.com/image.PNG", ".png"},
		{"https://example.com/image", ""},
		{"https://example.com/image.txt", ""}, // invalid extension
		{"https://example.com/", ""},
	}

	for _, tc := range testCases {
		result := downloader.GetExtensionFromURL(tc.url)
		if result != tc.expected {
			t.Errorf("For URL '%s', expected '%s', got '%s'", tc.url, tc.expected, result)
		}
	}
}

// TestCSVProcessing tests CSV file processing
func TestCSVProcessing(t *testing.T) {
	// Create a test CSV file
	csvData := `id,name,image_url,description
1,Test Image 1,https://example.com/image1.jpg,First test image
2,Test Image 2,https://example.com/image2.png,Second test image
3,Test Image 3,,Third test image with empty URL
4,Test Image 4,https://example.com/image4.gif,Fourth test image`

	testCSVFile := "test_data.csv"
	err := os.WriteFile(testCSVFile, []byte(csvData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}
	defer os.Remove(testCSVFile)

	// Test CSV reading
	file, err := os.Open(testCSVFile)
	if err != nil {
		t.Fatalf("Failed to open test CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		t.Fatalf("Failed to read CSV header: %v", err)
	}

	expectedHeader := []string{"id", "name", "image_url", "description"}
	if len(header) != len(expectedHeader) {
		t.Errorf("Expected header length %d, got %d", len(expectedHeader), len(header))
	}

	// Read rows
	rowCount := 0
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		rowCount++

		// Test that we have enough columns
		if len(row) < 3 {
			t.Errorf("Row %d: Expected at least 3 columns, got %d", rowCount, len(row))
		}

		// Test URL extraction (column 3, index 2)
		if len(row) >= 3 {
			url := strings.TrimSpace(row[2])
			if rowCount == 3 && url != "" {
				t.Errorf("Row %d: Expected empty URL, got '%s'", rowCount, url)
			}
		}
	}

	if rowCount != 4 {
		t.Errorf("Expected 4 rows, got %d", rowCount)
	}
}

// TestCommandLineArguments tests argument validation
func TestCommandLineArguments(t *testing.T) {
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test missing arguments
	os.Args = []string{"go-get-imgs"}
	// This would normally call main() and exit, so we test the validation logic directly

	// Test invalid column index
	invalidIndex := "abc"
	_, err := strconv.Atoi(invalidIndex)
	if err == nil {
		t.Error("Expected error for invalid column index 'abc', got nil")
	}

	// Test valid column index
	validIndex := "3"
	index, err := strconv.Atoi(validIndex)
	if err != nil {
		t.Errorf("Expected no error for valid column index '3', got %v", err)
	}
	if index != 3 {
		t.Errorf("Expected index 3, got %d", index)
	}
}

// TestFileOperations tests file and directory operations
func TestFileOperations(t *testing.T) {
	testDir := "test_operations"

	// Test directory creation
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Errorf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Test file creation
	testFile := filepath.Join(testDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}

	// Test file existence check
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Test file should exist")
	}

	// Test non-existent file check
	if _, err := os.Stat("non_existent_file.txt"); !os.IsNotExist(err) {
		t.Error("Non-existent file should return IsNotExist error")
	}
}

// BenchmarkDownloadImage benchmarks the download function
func BenchmarkDownloadImage(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte("benchmark image data"))
	}))
	defer server.Close()

	d := downloader.NewDownloader(30 * time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.DownloadImage(server.URL, "test_downloads", i)
	}
}

// TestIsValidURL tests the URL validation function
func TestIsValidURL(t *testing.T) {
	testCases := []struct {
		url      string
		expected bool
	}{
		{"https://example.com/image.jpg", true},
		{"http://example.com/image.png", true},
		{"ftp://example.com/image.gif", true},
		{"file:///path/to/image.jpg", true},
		{"invalid-url", false},
		{"", false},
		{"   ", false},
		{"  https://example.com/image.jpg  ", true},
		{"  invalid-url  ", false},
		{"https://", false}, // Incomplete URL
		{"http://", false},  // Incomplete URL
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("URL_%s", tc.url), func(t *testing.T) {
			result := utils.IsValidURL(tc.url)
			if result != tc.expected {
				t.Errorf("For URL '%s', expected %v, got %v", tc.url, tc.expected, result)
			}
		})
	}
}
