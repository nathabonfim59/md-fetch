package converter

import (
	"testing"
)

func TestConvertToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "simple paragraph",
			html:     "<p>Hello World</p>",
			expected: "Hello World",
		},
		{
			name:     "heading",
			html:     "<h1>Title</h1>",
			expected: "# Title",
		},
		{
			name:     "multiple elements",
			html:     "<h1>Title</h1><p>Content</p>",
			expected: "# Title\n\nContent",
		},
		{
			name:     "nested elements",
			html:     "<div><h1>Title</h1><p>Content with <strong>bold</strong> text</p></div>",
			expected: "# Title\n\nContent with **bold** text",
		},
		{
			name:     "links",
			html:     `<p>Visit <a href="https://example.com">Example</a></p>`,
			expected: "Visit [Example](https://example.com)",
		},
		{
			name:     "lists",
			html:     "<ul><li>Item 1</li><li>Item 2</li></ul>",
			expected: "- Item 1\n- Item 2",
		},
		{
			name:     "code blocks",
			html:     "<pre><code>func main() {\n    fmt.Println(\"Hello\")\n}</code></pre>",
			expected: "```\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
		},
		{
			name:     "empty input",
			html:     "",
			expected: "",
		},
		{
			name:     "invalid html",
			html:     "<p>Unclosed paragraph",
			expected: "Unclosed paragraph",  // Package tries to recover and outputs cleaned text
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToMarkdown([]byte(tt.html))
			if result != tt.expected {
				t.Errorf("\nexpected:\n%q\ngot:\n%q", tt.expected, result)
			}
		})
	}
}
