package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ContentType int

const (
	Html ContentType = iota
	Plaintext
	Json
)

type TagHandler interface {
	ShouldHandle(tag string) bool
	HandleTagStart(tag string, attrs map[string]string, writer *MarkdownWriter) bool
	HandleTagEnd(tag string, writer *MarkdownWriter)
}

type MarkdownWriter struct {
	buffer strings.Builder
}

func (w *MarkdownWriter) Write(text string) {
	w.buffer.WriteString(text)
}

func (w *MarkdownWriter) String() string {
	return w.buffer.String()
}

type WebpageChromeRemover struct{}

func (h *WebpageChromeRemover) ShouldHandle(tag string) bool {
	return tag == "head" || tag == "script" || tag == "style" || tag == "nav"
}

func (h *WebpageChromeRemover) HandleTagStart(tag string, attrs map[string]string, writer *MarkdownWriter) bool {
	return false // Skip content for these tags
}

func (h *WebpageChromeRemover) HandleTagEnd(tag string, writer *MarkdownWriter) {}

type ParagraphHandler struct{}

func (h *ParagraphHandler) ShouldHandle(tag string) bool {
	return tag == "p"
}

func (h *ParagraphHandler) HandleTagStart(tag string, attrs map[string]string, writer *MarkdownWriter) bool {
	writer.Write("\n\n")
	return true
}

func (h *ParagraphHandler) HandleTagEnd(tag string, writer *MarkdownWriter) {
	writer.Write("\n")
}

type HeadingHandler struct{}

func (h *HeadingHandler) ShouldHandle(tag string) bool {
	return strings.HasPrefix(tag, "h") && len(tag) == 2 && tag[1] >= '1' && tag[1] <= '6'
}

func (h *HeadingHandler) HandleTagStart(tag string, attrs map[string]string, writer *MarkdownWriter) bool {
	level := int(tag[1] - '0')
	writer.Write("\n" + strings.Repeat("#", level) + " ")
	return true
}

func (h *HeadingHandler) HandleTagEnd(tag string, writer *MarkdownWriter) {
	writer.Write("\n")
}

func convertHtmlToMarkdown(html []byte) string {
	handlers := []TagHandler{
		&WebpageChromeRemover{},
		&ParagraphHandler{},
		&HeadingHandler{},
	}

	writer := &MarkdownWriter{}

	// Simple HTML parsing and conversion
	content := string(html)
	var currentText strings.Builder

	for i := 0; i < len(content); {
		if content[i] == '<' {
			// Handle accumulated text
			if currentText.Len() > 0 {
				writer.Write(currentText.String())
				currentText.Reset()
			}

			// Find tag end
			tagEnd := strings.Index(content[i:], ">")
			if tagEnd == -1 {
				break
			}

			tagContent := content[i+1 : i+tagEnd]
			isClosing := strings.HasPrefix(tagContent, "/")
			if isClosing {
				tagContent = tagContent[1:]
			}

			// Extract tag name and attributes
			parts := strings.Fields(tagContent)
			tagName := strings.ToLower(parts[0])

			attrs := make(map[string]string)
			for _, attr := range parts[1:] {
				kv := strings.SplitN(attr, "=", 2)
				if len(kv) == 2 {
					attrs[kv[0]] = strings.Trim(kv[1], `"'`)
				}
			}

			if !isClosing {
				// Handle tag start
				for _, handler := range handlers {
					if handler.ShouldHandle(tagName) {
						handler.HandleTagStart(tagName, attrs, writer)
						break
					}
				}
			} else {
				// Handle tag end
				for _, handler := range handlers {
					if handler.ShouldHandle(tagName) {
						handler.HandleTagEnd(tagName, writer)
						break
					}
				}
			}

			i += tagEnd + 1
		} else {
			currentText.WriteByte(content[i])
			i++
		}
	}

	// Handle any remaining text
	if currentText.Len() > 0 {
		writer.Write(currentText.String())
	}

	return writer.String()
}

func fetchContent(url string) (string, error) {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("status error %d, response: %s", resp.StatusCode, string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return "", fmt.Errorf("missing Content-Type header")
	}

	var ct ContentType
	switch {
	case strings.Contains(contentType, "text/html"):
		ct = Html
	case strings.Contains(contentType, "text/plain"):
		ct = Plaintext
	case strings.Contains(contentType, "application/json"):
		ct = Json
	default:
		ct = Html
	}

	switch ct {
	case Html:
		return convertHtmlToMarkdown(body), nil
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

func main() {
	url := "https://example.com"
	content, err := fetchContent(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Fetched Content:")
	fmt.Println(content)
}
