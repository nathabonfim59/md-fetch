package fetcher

import (
	"bytes"
	"encoding/json"
	"fmt"
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
func FetchContent(url string, browserType string) (string, error) {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	b, err := browser.NewBrowser(browserType)
	if err != nil {
		return "", fmt.Errorf("error creating browser: %v", err)
	}

	body, err := b.Fetch(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL with %s: %v", b.Name(), err)
	}

	// Try to determine content type from first few bytes
	contentType := detectContentType(body)

	switch contentType {
	case Html:
		return converter.ConvertToMarkdown(body), nil
	case Plaintext:
		return string(body), nil
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
