package fetcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/nathabonfim59/md-fetch/internal/browser"
	"github.com/nathabonfim59/md-fetch/internal/converter"
)

type ContentType int

const (
	Html ContentType = iota
	Plaintext
	Json
)

// FetchContent retrieves and processes content from a URL using the specified browser
func FetchContent(urlStr string, browserType string) (string, error) {
	// Validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https" && !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://")) {
		// Try adding https:// prefix if no scheme is present
		if !strings.Contains(urlStr, "://") {
			urlStr = "https://" + urlStr
			parsedURL, err = url.Parse(urlStr)
			if err != nil {
				return "", fmt.Errorf("invalid URL %q: %v", urlStr, err)
			}
		} else {
			return "", fmt.Errorf("invalid URL %q: must use http or https scheme", urlStr)
		}
	}

	// Get browser instance
	var browserErr error
	var b browser.Browser
	if browserType == "" {
		b, browserErr = browser.GetDefaultBrowser()
	} else {
		b, browserErr = browser.NewBrowser(browserType)
	}
	if browserErr != nil {
		return "", fmt.Errorf("failed to initialize browser: %v", browserErr)
	}

	// Fetch content
	body, fetchErr := b.Fetch(urlStr)
	if fetchErr != nil {
		return "", fmt.Errorf("failed to fetch content: %v", fetchErr)
	}

	// Check for Chrome error pages
	bodyStr := string(body)
	if strings.Contains(bodyStr, "This site can't be reached") ||
		strings.Contains(bodyStr, "DNS_PROBE_FINISHED_NXDOMAIN") {
		return "", fmt.Errorf("failed to fetch content: site cannot be reached")
	}

	// Try to determine content type from first few bytes
	contentType := detectContentType(body)

	switch contentType {
	case Html:
		return converter.ConvertToMarkdown(body), nil
	case Plaintext:
		return bodyStr, nil
	case Json:
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			return "", fmt.Errorf("error formatting JSON: %v", err)
		}
		return "```json\n" + prettyJSON.String() + "\n```", nil
	default:
		return "", fmt.Errorf("unsupported content type")
	}
}

func detectContentType(content []byte) ContentType {
	// Simple content type detection based on content
	s := strings.TrimSpace(string(content))

	if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
		return Json
	} else if strings.Contains(s, "<html") || strings.Contains(s, "<!DOCTYPE html") {
		return Html
	}

	return Plaintext
}
