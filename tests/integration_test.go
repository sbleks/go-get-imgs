package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sbleks/go-get-imgs/internal/downloader"
	"github.com/sbleks/go-get-imgs/internal/utils"
)

// TestIntegrationCompleteWorkflow tests the complete application workflow
func TestIntegrationCompleteWorkflow(t *testing.T) {
	// Create test server that serves different images
	imageCounter := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageCounter++
		switch imageCounter {
		case 1:
			w.Header().Set("Content-Type", "image/jpeg")
		case 2:
			w.Header().Set("Content-Type", "image/png")
		case 3:
			w.Header().Set("Content-Type", "image/gif")
		default:
			w.Header().Set("Content-Type", "image/webp")
		}
		if _, err := w.Write([]byte(fmt.Sprintf("fake image data %d", imageCounter))); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create test CSV file
	csvData := fmt.Sprintf(`id,name,image_url,description
1,Test Image 1,%s/image1,First test image
2,Test Image 2,%s/image2,Second test image
3,Test Image 3,%s/image3,Third test image
4,Test Image 4,%s/image4,Fourth test image`, server.URL, server.URL, server.URL, server.URL)

	testCSVFile := "integration_test.csv"
	err := os.WriteFile(testCSVFile, []byte(csvData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}
	defer os.Remove(testCSVFile)

	// Create test downloads directory
	testDownloadsDir := "integration_downloads"
	err = os.MkdirAll(testDownloadsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test downloads directory: %v", err)
	}
	defer os.RemoveAll(testDownloadsDir)

	// Process CSV file manually (simulating main function logic)
	file, err := os.Open(testCSVFile)
	if err != nil {
		t.Fatalf("Failed to open test CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		t.Fatalf("Failed to read header: %v", err)
	}

	// Process each row
	successCount := 0
	errorCount := 0
	rowNum := 1

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		if len(row) < 3 {
			t.Errorf("Row %d: Not enough columns (need at least 3, got %d)", rowNum, len(row))
			errorCount++
			rowNum++
			continue
		}

		imageURL := strings.TrimSpace(row[2])
		if imageURL == "" {
			t.Errorf("Row %d: Empty URL in column 3", rowNum)
			errorCount++
			rowNum++
			continue
		}

		d := downloader.NewDownloader(30 * time.Second)
		if err := d.DownloadImage(imageURL, testDownloadsDir, rowNum); err != nil {
			t.Errorf("Row %d: Failed to download %s - %v", rowNum, imageURL, err)
			errorCount++
		} else {
			successCount++
		}

		rowNum++
	}

	// Verify results
	if successCount != 4 {
		t.Errorf("Expected 4 successful downloads, got %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}

	// Check that files were created with correct extensions
	expectedFiles := []string{
		"image_1.jpg",
		"image_2.png",
		"image_3.gif",
		"image_4.webp",
	}

	for _, expectedFile := range expectedFiles {
		filePath := filepath.Join(testDownloadsDir, expectedFile)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to be created", filePath)
		}
	}
}

// TestIntegrationWithErrors tests the application with various error scenarios
func TestIntegrationWithErrors(t *testing.T) {
	// Create test server that sometimes fails
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		switch requestCount {
		case 1:
			// Success
			w.Header().Set("Content-Type", "image/jpeg")
			if _, err := w.Write([]byte("success image")); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		case 2:
			// 404 error
			w.WriteHeader(http.StatusNotFound)
		case 3:
			// Success
			w.Header().Set("Content-Type", "image/png")
			if _, err := w.Write([]byte("success image")); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		case 4:
			// 500 error
			w.WriteHeader(http.StatusInternalServerError)
		default:
			// Success
			w.Header().Set("Content-Type", "image/gif")
			if _, err := w.Write([]byte("success image")); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		}
	}))
	defer server.Close()

	// Create test CSV file with some invalid URLs
	csvData := fmt.Sprintf(`id,name,image_url,description
1,Test Image 1,%s/image1,Should succeed
2,Test Image 2,%s/image2,Should fail with 404
3,Test Image 3,%s/image3,Should succeed
4,Test Image 4,%s/image4,Should fail with 500
5,Test Image 5,invalid-url,Should fail with invalid URL
6,Test Image 6,,Should fail with empty URL
7,Test Image 7,%s/image7,Should succeed`, server.URL, server.URL, server.URL, server.URL, server.URL)

	testCSVFile := "integration_errors_test.csv"
	err := os.WriteFile(testCSVFile, []byte(csvData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}
	defer os.Remove(testCSVFile)

	// Create test downloads directory
	testDownloadsDir := "integration_errors_downloads"
	err = os.MkdirAll(testDownloadsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test downloads directory: %v", err)
	}
	defer os.RemoveAll(testDownloadsDir)

	// Process CSV file
	file, err := os.Open(testCSVFile)
	if err != nil {
		t.Fatalf("Failed to open test CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	if _, err := reader.Read(); err != nil {
		t.Fatalf("Failed to read header: %v", err)
	}

	successCount := 0
	errorCount := 0
	rowNum := 1

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		if len(row) < 3 {
			errorCount++
			rowNum++
			continue
		}

		imageURL := strings.TrimSpace(row[2])
		if imageURL == "" {
			errorCount++
			rowNum++
			continue
		}

		// Validate URL format (matching main application logic)
		if !utils.IsValidURL(imageURL) {
			errorCount++
			rowNum++
			continue
		}

		d := downloader.NewDownloader(30 * time.Second)
		if err := d.DownloadImage(imageURL, testDownloadsDir, rowNum); err != nil {
			errorCount++
		} else {
			successCount++
		}

		rowNum++
	}

	// Verify results - should have 3 successes and 4 errors
	if successCount != 3 {
		t.Errorf("Expected 3 successful downloads, got %d", successCount)
	}

	if errorCount != 4 {
		t.Errorf("Expected 4 errors, got %d", errorCount)
	}

	// Check that only successful files were created
	expectedFiles := []string{
		"image_1.jpg",
		"image_3.png",
		"image_7.gif",
	}

	for _, expectedFile := range expectedFiles {
		filePath := filepath.Join(testDownloadsDir, expectedFile)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to be created", filePath)
		}
	}

	// Check that failed files were not created
	unexpectedFiles := []string{
		"image_2.jpg",
		"image_4.jpg",
		"image_5.jpg",
		"image_6.jpg",
	}

	for _, unexpectedFile := range unexpectedFiles {
		filePath := filepath.Join(testDownloadsDir, unexpectedFile)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("Expected file %s to NOT be created", filePath)
		}
	}
}

// TestIntegrationLargeFile tests with a larger CSV file
func TestIntegrationLargeFile(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		if _, err := w.Write([]byte("large file test image")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create a larger CSV file (100 rows)
	var csvBuilder strings.Builder
	csvBuilder.WriteString("id,name,image_url,description\n")

	for i := 1; i <= 100; i++ {
		csvBuilder.WriteString(fmt.Sprintf("%d,Image %d,%s/image%d,Description %d\n",
			i, i, server.URL, i, i))
	}

	testCSVFile := "large_integration_test.csv"
	err := os.WriteFile(testCSVFile, []byte(csvBuilder.String()), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test CSV file: %v", err)
	}
	defer os.Remove(testCSVFile)

	// Create test downloads directory
	testDownloadsDir := "large_integration_downloads"
	err = os.MkdirAll(testDownloadsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test downloads directory: %v", err)
	}
	defer os.RemoveAll(testDownloadsDir)

	// Process CSV file
	file, err := os.Open(testCSVFile)
	if err != nil {
		t.Fatalf("Failed to open test CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	if _, err := reader.Read(); err != nil {
		t.Fatalf("Failed to read header: %v", err)
	}

	successCount := 0
	errorCount := 0
	rowNum := 1

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		if len(row) < 3 {
			errorCount++
			rowNum++
			continue
		}

		imageURL := strings.TrimSpace(row[2])
		if imageURL == "" {
			errorCount++
			rowNum++
			continue
		}

		// Validate URL format (matching main application logic)
		if !utils.IsValidURL(imageURL) {
			errorCount++
			rowNum++
			continue
		}

		d := downloader.NewDownloader(30 * time.Second)
		if err := d.DownloadImage(imageURL, testDownloadsDir, rowNum); err != nil {
			errorCount++
		} else {
			successCount++
		}

		rowNum++
	}

	// Verify results
	if successCount != 100 {
		t.Errorf("Expected 100 successful downloads, got %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}

	// Check that all files were created
	for i := 1; i <= 100; i++ {
		expectedFile := fmt.Sprintf("image_%d.jpg", i)
		filePath := filepath.Join(testDownloadsDir, expectedFile)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to be created", filePath)
		}
	}
}
