package browser

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// DefaultCleaningOptions returns the default cleaning configuration
func DefaultCleaningOptions() *CleaningOptions {
	return &CleaningOptions{
		KeepHeader:    false,
		KeepFooter:    false,
		KeepNav:       false,
		KeepStyles:    false,
		KeepComments:  false,
	}
}

// StripJavaScript removes JavaScript code from HTML content
func StripJavaScript(content []byte) []byte {
	return CleanHTML(content, DefaultCleaningOptions())
}

// CleanHTML removes unwanted elements from HTML content based on options
func CleanHTML(content []byte, opts *CleaningOptions) []byte {
	// Convert content to string for easier regex operations
	contentStr := string(content)

	// Remove script tags and their content
	scriptPattern := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	contentStr = scriptPattern.ReplaceAllString(contentStr, "")

	// Remove style tags if not keeping styles
	if !opts.KeepStyles {
		stylePattern := regexp.MustCompile(`(?is)<style.*?>.*?</style>`)
		contentStr = stylePattern.ReplaceAllString(contentStr, "")
	}

	// Remove comments if not keeping them
	if !opts.KeepComments {
		commentPattern := regexp.MustCompile(`(?is)<!--.*?-->`)
		contentStr = commentPattern.ReplaceAllString(contentStr, "")
	}

	// Remove inline JavaScript attributes
	jsAttrPattern := regexp.MustCompile(`(?i)(on\w+)="[^"]*"`)
	contentStr = jsAttrPattern.ReplaceAllString(contentStr, "")

	// Remove various JavaScript patterns
	jsPatterns := []*regexp.Regexp{
		// Google-specific patterns
		regexp.MustCompile(`(?s)var\s+_g\s*=\s*\{\s*kEI\s*:[^}]*\}\s*;`),
		regexp.MustCompile(`(?s)var\s+google\s*=\s*\{\s*[^}]*\}\s*;`),
		regexp.MustCompile(`(?s)google\.[a-zA-Z_$][0-9a-zA-Z_$]*\s*=\s*[^;]*;`),
		regexp.MustCompile(`(?s)\(\s*function\s*\(\)\s*\{\s*var\s+a\s*=\s*window\.innerWidth[^}]*\}\s*\)\s*\(\)\s*;`),
		regexp.MustCompile(`(?s)window\._cshid\s*&&[^;]*;`),
		
		// Anonymous function calls with window assignments
		regexp.MustCompile(`(?s)\(\s*function\s*\(\)\s*\{\s*var\s+[a-zA-Z_$][0-9a-zA-Z_$]*\s*=\s*\{[^}]*\}\s*;[^}]*\}\s*\)\s*\(\s*\)\s*;`),
		
		// RLQ function calls
		regexp.MustCompile(`(?s)\(\s*RLQ\s*=\s*window\.RLQ\s*\|\|\s*\[\]\s*\)\.push\s*\(\s*function\s*\([^)]*\)\s*\{[^}]*\}\s*\)\s*;`),
		
		// Self-executing functions with call
		regexp.MustCompile(`(?s)\(\s*function\s*\([^)]*\)\s*\{[^}]*\}\s*\)\s*\.\s*call\s*\([^)]*\)\s*;`),
		
		// Self-executing function expressions
		regexp.MustCompile(`(?s)\(\s*function\s*\([^)]*\)\s*\{[^}]*\}\s*\)\s*\([^)]*\)\s*;?`),
		
		// Window property assignments
		regexp.MustCompile(`(?m)^[ \t]*window\.[a-zA-Z_$][0-9a-zA-Z_$]*\s*=\s*[^;]*;`),
		
		// Document event listeners
		regexp.MustCompile(`(?s)document\.[a-zA-Z_$][0-9a-zA-Z_$]*\.addEventListener\s*\([^)]*\)\s*;`),
		
		// Inline JSON-LD
		regexp.MustCompile(`(?s)\{\s*"@context"\s*:\s*"https?:\\?/\\?/schema\.org"[^}]*\}`),
		
		// CSS definitions and styles
		regexp.MustCompile(`(?m)^[ \t]*#[a-zA-Z][0-9a-zA-Z_-]*\s*\{[^}]*\}`),
		regexp.MustCompile(`(?m)^[ \t]*\.[a-zA-Z][0-9a-zA-Z_-]*\s*\{[^}]*\}`),
		regexp.MustCompile(`(?m)^[ \t]*[a-zA-Z][0-9a-zA-Z_-]*\s*\{[^}]*\}`),
		regexp.MustCompile(`(?m)^[ \t]*@media[^{]*\{[^}]*\}`),
		
		// Variable declarations
		regexp.MustCompile(`(?m)^[ \t]*var\s+[a-zA-Z_$][0-9a-zA-Z_$]*\s*=[^;]*;`),
		
		// Function declarations and expressions
		regexp.MustCompile(`(?m)^[ \t]*function\s+[a-zA-Z_$][0-9a-zA-Z_$]*\s*\([^)]*\)\s*\{[^}]*\}`),
		regexp.MustCompile(`\([^)]*\)\s*=>\s*\{[^}]*\}`),
		regexp.MustCompile(`\(\s*function\s*\([^)]*\)\s*\{[^}]*\}\s*\)`),
		
		// Comments
		regexp.MustCompile(`(?m)^[ \t]*//.*$`),
		regexp.MustCompile(`(?s)/\*.*?\*/`),
	}

	for _, pattern := range jsPatterns {
		contentStr = pattern.ReplaceAllString(contentStr, "")
	}

	// Parse and clean the HTML
	doc, err := html.Parse(strings.NewReader(contentStr))
	if err != nil {
		return content // Return original content if parsing fails
	}

	var buf bytes.Buffer
	cleanNode(&buf, doc, opts)
	return buf.Bytes()
}

func cleanNode(w io.Writer, n *html.Node, opts *CleaningOptions) bool {
	if n.Type == html.ElementNode {
		// Check if we should skip this node based on options
		if shouldSkipNode(n, opts) {
			return false
		}

		// Clean inline styles if not keeping them
		if !opts.KeepStyles {
			removeStyleAttr(n)
		}

		// Clean JavaScript URLs
		if n.Data == "a" {
			cleanJavaScriptURLs(n)
		}
	}

	// Write opening tag
	if n.Type == html.ElementNode {
		io.WriteString(w, "<"+n.Data)
		for _, attr := range n.Attr {
			io.WriteString(w, " "+attr.Key+"=\""+html.EscapeString(attr.Val)+"\"")
		}
		if voidElements[n.Data] {
			io.WriteString(w, "/>")
			return true
		}
		io.WriteString(w, ">")
	} else if n.Type == html.TextNode {
		// Clean any remaining JavaScript-like content in text nodes
		text := cleanTextContent(n.Data)
		io.WriteString(w, html.EscapeString(text))
	} else if n.Type == html.CommentNode && opts.KeepComments {
		io.WriteString(w, "<!--"+n.Data+"-->")
	}

	// Process children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if cleanNode(w, c, opts) {
			continue
		}
	}

	// Write closing tag
	if n.Type == html.ElementNode && !voidElements[n.Data] {
		io.WriteString(w, "</"+n.Data+">")
	}

	return true
}

func cleanTextContent(text string) string {
	// Remove any JavaScript-like content from text
	patterns := []*regexp.Regexp{
		// RLQ function calls
		regexp.MustCompile(`\(\s*RLQ\s*=\s*window\.RLQ\s*\|\|\s*\[\]\s*\)\.push\s*\(\s*function\s*\([^)]*\)\s*\{[^}]*\}\s*\)\s*;`),
		
		// JSON-LD data
		regexp.MustCompile(`\{\s*"@context"\s*:\s*"https?:\\?/\\?/schema\.org"[^}]*\}`),
		
		// Self-executing functions with call
		regexp.MustCompile(`\(\s*function\s*\([^)]*\)\s*\{[^}]*\}\s*\)\s*\.\s*call\s*\([^)]*\)\s*;`),
		
		// Self-executing function expressions
		regexp.MustCompile(`\(\s*function\s*\([^)]*\)\s*\{[^}]*\}\s*\)\s*\([^)]*\)\s*;?`),
		
		// Variable declarations
		regexp.MustCompile(`var\s+[a-zA-Z_$][0-9a-zA-Z_$]*\s*=[^;]*;`),
		
		// Function declarations
		regexp.MustCompile(`function\s+[a-zA-Z_$][0-9a-zA-Z_$]*\s*\([^)]*\)\s*\{[^}]*\}`),
		
		// Window assignments
		regexp.MustCompile(`window\.[a-zA-Z_$][0-9a-zA-Z_$]*\s*=\s*[^;]*;`),
		
		// Document event listeners
		regexp.MustCompile(`document\.[a-zA-Z_$][0-9a-zA-Z_$]*\.addEventListener\s*\([^)]*\)\s*;`),
		
		// Inline function expressions
		regexp.MustCompile(`\([^)]*\)\s*=>\s*\{[^}]*\}`),
		regexp.MustCompile(`\(\s*function\s*\([^)]*\)\s*\{[^}]*\}\s*\)`),
	}

	for _, pattern := range patterns {
		text = pattern.ReplaceAllString(text, "")
	}

	return text
}

func shouldSkipNode(n *html.Node, opts *CleaningOptions) bool {
	switch n.Data {
	case "header":
		return !opts.KeepHeader
	case "footer":
		return !opts.KeepFooter
	case "nav":
		return !opts.KeepNav
	case "style":
		return !opts.KeepStyles
	}
	return false
}

func removeStyleAttr(n *html.Node) {
	for i := 0; i < len(n.Attr); i++ {
		if n.Attr[i].Key == "style" {
			// Remove the style attribute
			n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			i--
		}
	}
}

func cleanJavaScriptURLs(n *html.Node) {
	for i := 0; i < len(n.Attr); i++ {
		if n.Attr[i].Key == "href" && strings.HasPrefix(strings.TrimSpace(strings.ToLower(n.Attr[i].Val)), "javascript:") {
			n.Attr[i].Val = "#"
		}
	}
}

// List of void elements that don't need closing tags
var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}
