package browser

import (
	"strings"
	"testing"
)

func TestStripJavaScript(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		excludes []string
	}{
		{
			name: "script tags",
			input: `<html>
				<head><script>alert('test');</script></head>
				<body>Content</body>
			</html>`,
			contains: []string{"Content"},
			excludes: []string{"<script", "alert"},
		},
		{
			name: "event handlers",
			input: `<button onclick="alert('click')" onmouseover="hover()">Click me</button>`,
			contains: []string{"<button", "Click me", "</button>"},
			excludes: []string{"onclick", "onmouseover", "alert", "hover"},
		},
		{
			name: "javascript urls",
			input: `<a href="javascript:void(0)" onclick="click()">Link</a>`,
			contains: []string{`<a href="#"`, "Link", "</a>"},
			excludes: []string{"javascript:", "onclick", "click()"},
		},
		{
			name: "mixed content",
			input: `<div>
				<script>var x = 1;</script>
				<p>Text</p>
				<button onclick="click()">Button</button>
				<script type="text/javascript">alert('end');</script>
			</div>`,
			contains: []string{"<div", "<p>Text</p>", "<button", "Button", "</button>", "</div>"},
			excludes: []string{"<script", "var x = 1", "onclick", "click()", "alert"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(StripJavaScript([]byte(tt.input)))
			
			// Check for required content
			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("Expected content %q not found in result:\n%s", s, result)
				}
			}

			// Check for excluded content
			for _, s := range tt.excludes {
				if strings.Contains(result, s) {
					t.Errorf("Unexpected content %q found in result:\n%s", s, result)
				}
			}
		})
	}
}

func TestJavaScriptHandling(t *testing.T) {
	// Test with a simple JavaScript-enabled page
	html := `
		<html>
			<head>
				<script>document.write('Dynamic content');</script>
			</head>
			<body>
				<div id="content">Static content</div>
			</body>
		</html>`

	content := StripJavaScript([]byte(html))
	if strings.Contains(string(content), "<script") {
		t.Error("JavaScript tags found in stripped content")
	}
	if !strings.Contains(string(content), "Static content") {
		t.Error("Static content missing from stripped content")
	}
}
