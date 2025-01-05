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
	// Remove script tags and their content
	scriptPattern := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	content = scriptPattern.ReplaceAll(content, []byte(""))

	// Remove style tags if not keeping styles
	if !opts.KeepStyles {
		stylePattern := regexp.MustCompile(`(?is)<style.*?>.*?</style>`)
		content = stylePattern.ReplaceAll(content, []byte(""))
	}

	// Remove comments if not keeping them
	if !opts.KeepComments {
		commentPattern := regexp.MustCompile(`(?is)<!--.*?-->`)
		content = commentPattern.ReplaceAll(content, []byte(""))
	}

	// Remove inline JavaScript attributes
	jsAttrPattern := regexp.MustCompile(`(?i)(on\w+)="[^"]*"`)
	content = jsAttrPattern.ReplaceAll(content, []byte(""))

	// Parse and clean the HTML
	doc, err := html.Parse(bytes.NewReader(content))
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
		io.WriteString(w, html.EscapeString(n.Data))
	} else if n.Type == html.CommentNode && opts.KeepComments {
		io.WriteString(w, "<!--"+n.Data+"-->")
	}

	// Process children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if cleanNode(w, c, opts) {
			// If child was processed, move to next sibling
			continue
		}
	}

	// Write closing tag
	if n.Type == html.ElementNode && !voidElements[n.Data] {
		io.WriteString(w, "</"+n.Data+">")
	}

	return true
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
