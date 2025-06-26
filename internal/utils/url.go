package utils

import "strings"

// IsValidURL checks if a string is a valid URL with a supported scheme
func IsValidURL(url string) bool {
	trimmed := strings.TrimSpace(url)
	if trimmed == "" {
		return false
	}

	// Check for supported URL schemes
	validSchemes := []string{"http://", "https://", "ftp://", "file://"}
	hasValidScheme := false
	for _, scheme := range validSchemes {
		if strings.HasPrefix(trimmed, scheme) {
			hasValidScheme = true
			break
		}
	}

	if !hasValidScheme {
		return false
	}

	// For file:// URLs, check if there's a path after the scheme
	if strings.HasPrefix(trimmed, "file://") {
		// file:// URLs should have at least a path (e.g., file:///path)
		return len(trimmed) > 7
	}

	// For other schemes, check if there's a domain after the scheme
	// Remove the scheme and check if there's content after it
	for _, scheme := range validSchemes {
		if strings.HasPrefix(trimmed, scheme) {
			afterScheme := trimmed[len(scheme):]
			// Should have at least a domain (e.g., example.com)
			return len(afterScheme) > 0 && !strings.Contains(afterScheme, " ")
		}
	}

	return false
}
