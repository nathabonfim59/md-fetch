package converter

import (
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// ConvertToMarkdown converts HTML content to Markdown format
func ConvertToMarkdown(html []byte) string {
	// Create a new converter with plugins
	conv, err := htmltomarkdown.ConvertString(string(html))

	if err != nil {
		return string(html)
	}

	return conv
}
