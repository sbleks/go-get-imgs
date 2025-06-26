package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	csvpkg "github.com/sbleks/go-get-imgs/internal/csv"
	"github.com/sbleks/go-get-imgs/internal/downloader"
)

// TestHelper provides common test utilities
type TestHelper struct {
	t            *testing.T
	server       *httptest.Server
	cleanupFuncs []func()
}

// NewTestHelper creates a new test helper instance
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t:            t,
		cleanupFuncs: make([]func(), 0),
	}
}

// CreateTestServer creates a test HTTP server
func (th *TestHelper) CreateTestServer(handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	th.server = server
	th.cleanupFuncs = append(th.cleanupFuncs, server.Close)
	return server
}

// CreateTestServerWithImages creates a test server that serves different image types
func (th *TestHelper) CreateTestServerWithImages() *httptest.Server {
	imageCounter := 0
	server := th.CreateTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageCounter++
		switch imageCounter {
		case 1:
			w.Header().Set("Content-Type", "image/jpeg")
		case 2:
			w.Header().Set("Content-Type", "image/png")
		case 3:
			w.Header().Set("Content-Type", "image/gif")
		case 4:
			w.Header().Set("Content-Type", "image/webp")
		default:
			w.Header().Set("Content-Type", "image/jpeg")
		}
		if _, err := w.Write([]byte(fmt.Sprintf("fake image data %d", imageCounter))); err != nil {
			th.t.Errorf("Failed to write response: %v", err)
		}
	}))
	return server
}

// CreateTestCSV creates a test CSV file with the given data
func (th *TestHelper) CreateTestCSV(filename string, data string) string {
	err := os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		th.t.Fatalf("Failed to create test CSV file: %v", err)
	}

	th.cleanupFuncs = append(th.cleanupFuncs, func() {
		os.Remove(filename)
	})

	return filename
}

// CreateTestCSVWithURLs creates a test CSV file with image URLs
func (th *TestHelper) CreateTestCSVWithURLs(filename string, baseURL string, count int) string {
	var csvBuilder strings.Builder
	csvBuilder.WriteString("id,name,image_url,description\n")

	for i := 1; i <= count; i++ {
		csvBuilder.WriteString(fmt.Sprintf("%d,Image %d,%s/image%d,Description %d\n",
			i, i, baseURL, i, i))
	}

	return th.CreateTestCSV(filename, csvBuilder.String())
}

// CreateTestDirectory creates a test directory
func (th *TestHelper) CreateTestDirectory(dirname string) string {
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		th.t.Fatalf("Failed to create test directory: %v", err)
	}

	th.cleanupFuncs = append(th.cleanupFuncs, func() {
		os.RemoveAll(dirname)
	})

	return dirname
}

// ProcessCSVFile processes a CSV file and returns success/error counts
func (th *TestHelper) ProcessCSVFile(csvFile string, downloadDir string, urlColumnIndex int) (int, int) {
	processor := csvpkg.NewProcessor()
	downloader := downloader.NewDownloader(30 * time.Second)

	result, err := processor.ProcessCSV(csvFile, urlColumnIndex, func(url string, rowNum int) error {
		return downloader.DownloadImage(url, downloadDir, rowNum)
	})

	if err != nil {
		th.t.Fatalf("Failed to process CSV file: %v", err)
	}

	return result.SuccessCount, result.ErrorCount
}

// AssertFileExists checks if a file exists
func (th *TestHelper) AssertFileExists(filepath string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		th.t.Errorf("Expected file %s to exist", filepath)
	}
}

// AssertFileNotExists checks if a file does not exist
func (th *TestHelper) AssertFileNotExists(filepath string) {
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		th.t.Errorf("Expected file %s to NOT exist", filepath)
	}
}

// AssertFilesExist checks if multiple files exist
func (th *TestHelper) AssertFilesExist(directory string, filenames []string) {
	for _, filename := range filenames {
		filepath := filepath.Join(directory, filename)
		th.AssertFileExists(filepath)
	}
}

// AssertFilesNotExist checks if multiple files do not exist
func (th *TestHelper) AssertFilesNotExist(directory string, filenames []string) {
	for _, filename := range filenames {
		filepath := filepath.Join(directory, filename)
		th.AssertFileNotExists(filepath)
	}
}

// Cleanup performs cleanup operations
func (th *TestHelper) Cleanup() {
	for _, cleanup := range th.cleanupFuncs {
		cleanup()
	}
}

// ValidateCSVStructure validates the structure of a CSV file
func (th *TestHelper) ValidateCSVStructure(filename string, expectedColumns int) {
	processor := csvpkg.NewProcessor()
	err := processor.ValidateCSVStructure(filename, expectedColumns)
	if err != nil {
		th.t.Errorf("CSV validation failed: %v", err)
	}
}

// CreateMockImageServer creates a server that serves mock images with specific behaviors
func (th *TestHelper) CreateMockImageServer(behaviors map[int]MockBehavior) *httptest.Server {
	requestCount := 0
	server := th.CreateTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		behavior, exists := behaviors[requestCount]
		if !exists {
			// Default behavior
			w.Header().Set("Content-Type", "image/jpeg")
			if _, err := w.Write([]byte("default image data")); err != nil {
				th.t.Errorf("Failed to write response: %v", err)
			}
			return
		}

		switch behavior.Type {
		case "success":
			w.Header().Set("Content-Type", behavior.ContentType)
			if _, err := w.Write([]byte(behavior.Data)); err != nil {
				th.t.Errorf("Failed to write response: %v", err)
			}
		case "error":
			w.WriteHeader(behavior.StatusCode)
		case "timeout":
			// Simulate timeout by not responding
			select {}
		}
	}))
	return server
}

// MockBehavior defines how a mock server should behave for a specific request
type MockBehavior struct {
	Type        string // "success", "error", "timeout"
	ContentType string
	Data        string
	StatusCode  int
}
