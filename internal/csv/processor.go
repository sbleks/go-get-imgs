package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Processor handles CSV file processing operations
type Processor struct{}

// NewProcessor creates a new CSV processor instance
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessResult contains the results of CSV processing
type ProcessResult struct {
	SuccessCount int
	ErrorCount   int
	TotalRows    int
}

// ProcessCSV processes a CSV file and returns processing results
func (p *Processor) ProcessCSV(csvFile string, urlColumnIndex int, downloadFunc func(url string, rowNum int) error) (*ProcessResult, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	if len(header) < urlColumnIndex {
		return nil, fmt.Errorf("expected at least %d columns in header, got %d", urlColumnIndex, len(header))
	}

	result := &ProcessResult{}
	rowNum := 1

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		result.TotalRows++

		if len(row) < urlColumnIndex {
			result.ErrorCount++
			rowNum++
			continue
		}

		imageURL := strings.TrimSpace(row[urlColumnIndex-1])
		if imageURL == "" {
			result.ErrorCount++
			rowNum++
			continue
		}

		if err := downloadFunc(imageURL, rowNum); err != nil {
			result.ErrorCount++
		} else {
			result.SuccessCount++
		}

		rowNum++
	}

	return result, nil
}

// ValidateCSVStructure validates the structure of a CSV file
func (p *Processor) ValidateCSVStructure(filename string, expectedColumns int) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open CSV file for validation: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %v", err)
	}

	if len(header) < expectedColumns {
		return fmt.Errorf("expected at least %d columns in header, got %d", expectedColumns, len(header))
	}

	// Count rows
	rowCount := 0
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		rowCount++

		if len(row) < expectedColumns {
			return fmt.Errorf("row %d: expected at least %d columns, got %d", rowCount, expectedColumns, len(row))
		}
	}

	return nil
}
