package browser

import (
	"strings"
	"testing"
)

func TestBrowserRedirect(t *testing.T) {
	// Get default browser
	b, err := GetDefaultBrowser()
	if err != nil {
		t.Skipf("No browser available for testing: %v", err)
		return
	}

	// Test redirect from http to https
	content, err := b.Fetch("http://github.com")
	if err != nil {
		t.Errorf("Failed to fetch from redirecting URL: %v", err)
		return
	}

	// Verify we got actual content
	if len(content) == 0 {
		t.Error("Got empty content from redirecting URL")
		return
	}

	// Verify it's HTML content (handle different formats)
	contentStr := strings.ToLower(string(content))
	if !strings.Contains(contentStr, "<!doctype html") &&
		!strings.Contains(contentStr, "<html") &&
		!strings.Contains(contentStr, "<body") {
		t.Error("Content does not appear to be HTML")
	}
}
