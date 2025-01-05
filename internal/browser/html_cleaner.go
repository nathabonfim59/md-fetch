package browser

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// StripJavaScript removes JavaScript code from HTML content
func StripJavaScript(content []byte) []byte {
	// Remove <script> tags and their content
	scriptPattern := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	content = scriptPattern.ReplaceAll(content, []byte(""))

	// Remove inline JavaScript attributes
	jsAttrPattern := regexp.MustCompile(`(?i)(on\w+)="[^"]*"`)
	content = jsAttrPattern.ReplaceAll(content, []byte(""))

	// Parse and clean the HTML
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return content // Return original content if parsing fails
	}

	var buf bytes.Buffer
	cleanNode(&buf, doc)
	return buf.Bytes()
}

func cleanNode(w io.Writer, n *html.Node) {
	if n.Type == html.ElementNode {
		// Skip JavaScript URLs
		if n.Data == "a" {
			for i := 0; i < len(n.Attr); i++ {
				if n.Attr[i].Key == "href" && strings.HasPrefix(strings.TrimSpace(strings.ToLower(n.Attr[i].Val)), "javascript:") {
					n.Attr[i].Val = "#"
				}
			}
		}
	}

	// Render the node
	if n.Type == html.ElementNode {
		io.WriteString(w, "<"+n.Data)
		for _, attr := range n.Attr {
			io.WriteString(w, " "+attr.Key+"=\""+attr.Val+"\"")
		}
		if voidElements[n.Data] {
			io.WriteString(w, "/>")
		} else {
			io.WriteString(w, ">")
		}
	} else if n.Type == html.TextNode {
		io.WriteString(w, html.EscapeString(n.Data))
	}

	// Process children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cleanNode(w, c)
	}

	// Close non-void elements
	if n.Type == html.ElementNode && !voidElements[n.Data] {
		io.WriteString(w, "</"+n.Data+">")
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
