package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Downloader handles image downloading operations
type Downloader struct {
	client *http.Client
}

// NewDownloader creates a new downloader instance
func NewDownloader(timeout time.Duration) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// DownloadImage downloads an image from a URL and saves it to the specified directory
func (d *Downloader) DownloadImage(url, downloadDir string, rowNum int) error {
	resp, err := d.client.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	extension := getExtensionFromContentType(contentType)
	if extension == "" {
		extension = GetExtensionFromURL(url)
		if extension == "" {
			extension = ".jpg"
		}
	}

	filename := fmt.Sprintf("image_%d%s", rowNum, extension)
	filepath := filepath.Join(downloadDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// getExtensionFromContentType determines file extension from HTTP content-type header
func getExtensionFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "image/jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "image/jpg"):
		return ".jpg"
	case strings.Contains(contentType, "image/png"):
		return ".png"
	case strings.Contains(contentType, "image/gif"):
		return ".gif"
	case strings.Contains(contentType, "image/webp"):
		return ".webp"
	case strings.Contains(contentType, "image/bmp"):
		return ".bmp"
	case strings.Contains(contentType, "image/tiff"):
		return ".tiff"
	default:
		return ""
	}
}

// GetExtensionFromContentType determines file extension from HTTP content-type header
func GetExtensionFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "image/jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "image/jpg"):
		return ".jpg"
	case strings.Contains(contentType, "image/png"):
		return ".png"
	case strings.Contains(contentType, "image/gif"):
		return ".gif"
	case strings.Contains(contentType, "image/webp"):
		return ".webp"
	case strings.Contains(contentType, "image/bmp"):
		return ".bmp"
	case strings.Contains(contentType, "image/tiff"):
		return ".tiff"
	default:
		return ""
	}
}

// GetExtensionFromURL determines file extension from URL path
func GetExtensionFromURL(url string) string {
	ext := filepath.Ext(url)
	if ext != "" {
		ext = strings.ToLower(ext)
		validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff", ".tif"}
		for _, validExt := range validExts {
			if ext == validExt {
				return ext
			}
		}
	}
	return ""
}
