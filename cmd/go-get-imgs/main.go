package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sbleks/go-get-imgs/internal/csv"
	"github.com/sbleks/go-get-imgs/internal/downloader"
	"github.com/sbleks/go-get-imgs/internal/utils"
)

// Version information for cross-platform builds
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go-get-imgs <csv-file> <url-column-index>")
		fmt.Println("Example: go-get-imgs data.csv 3")
		fmt.Printf("Version: %s (Built: %s, Commit: %s)\n", Version, BuildTime, GitCommit)
		os.Exit(1)
	}

	csvFile := os.Args[1]
	urlColumnIndex, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: Invalid URL column index: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		fmt.Printf("Error: CSV file '%s' does not exist\n", csvFile)
		os.Exit(1)
	}

	// Create downloads directory
	downloadsDir := "downloads"
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		fmt.Printf("Error creating downloads directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize components
	processor := csv.NewProcessor()
	downloader := downloader.NewDownloader(30 * time.Second)

	// Process CSV file
	result, err := processor.ProcessCSV(csvFile, urlColumnIndex, func(url string, rowNum int) error {
		// Validate URL format
		if !utils.IsValidURL(url) {
			return fmt.Errorf("invalid URL format: %s", url)
		}

		fmt.Printf("Downloading row %d: %s\n", rowNum, url)
		return downloader.DownloadImage(url, downloadsDir, rowNum)
	})

	if err != nil {
		fmt.Printf("Error processing CSV file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDownload Summary:\n")
	fmt.Printf("✅ Successful downloads: %d\n", result.SuccessCount)
	fmt.Printf("❌ Failed downloads: %d\n", result.ErrorCount)
	fmt.Printf("📁 Images saved to: %s/\n", downloadsDir)
}
