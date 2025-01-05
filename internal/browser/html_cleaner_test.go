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
		{
			name: "anonymous functions and json",
			input: `<html>
				<head>
					<script>
						(RLQ=window.RLQ||[]).push(function(){mw.config.set({"wgHostname":"test"});});
					</script>
				</head>
				<body>
					<div>Content</div>
					{"@context":"https://schema.org","@type":"Article","name":"Test"}
					(function(x) { console.log(x); })();
				</body>
			</html>`,
			contains: []string{"<div>Content</div>"},
			excludes: []string{
				"RLQ",
				"window.RLQ",
				"mw.config",
				"wgHostname",
				"schema.org",
				"@context",
				"function",
				"console.log",
			},
		},
		{
			name: "complex javascript patterns",
			input: `<html>
				<head>
					<script>
						(function(){var _g={kEI:'test'};(function(){window.google=_g;}).call(this);})();
						(RLQ=window.RLQ||[]).push(function(){mw.config.set({"test":"value"});});
						google.x=function(a,b){google.y[c]=[a,b];return!1};
						document.documentElement.addEventListener("submit",function(b){});
						var g=this||self;var k,l=(k=g.mei)!=null?k:1;
						#gbar,#guser{font-size:13px;padding-top:1px !important;}
						window.onerror=function(a,b,d,n,e){return null};
					</script>
				</head>
				<body>
					<div>Content</div>
					{"@context":"https://schema.org","@type":"Article","name":"Test"}
					<script>
						(function(){
							var src='/images/test.png';
							document.body.onload = function(){};
						})();
					</script>
				</body>
			</html>`,
			contains: []string{"<div>Content</div>"},
			excludes: []string{
				"function",
				"window.RLQ",
				"mw.config",
				"google.x",
				"addEventListener",
				"document.body",
				"onerror",
				"schema.org",
				"@context",
				"#gbar",
				"font-size",
				"var src",
				"onload",
			},
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

func TestCleanHTML(t *testing.T) {
	html := `
		<html>
			<head>
				<style>.header { color: red; }</style>
				<script>document.write('Dynamic');</script>
			</head>
			<body>
				<header>Site Header</header>
				<nav>Navigation Menu</nav>
				<div id="content" style="color: blue;">
					Main Content
					<!-- Comment -->
				</div>
				<footer>Site Footer</footer>
			</body>
		</html>`

	tests := []struct {
		name     string
		opts     *CleaningOptions
		contains []string
		excludes []string
	}{
		{
			name: "default options",
			opts: DefaultCleaningOptions(),
			contains: []string{"Main Content", "<div", "</div>"},
			excludes: []string{
				"<header", "Site Header",
				"<nav", "Navigation Menu",
				"<footer", "Site Footer",
				"<style", "color: red",
				"color: blue",
				"<!-- Comment -->",
				"<script",
			},
		},
		{
			name: "keep header only",
			opts: &CleaningOptions{
				KeepHeader: true,
			},
			contains: []string{
				"<header>Site Header</header>",
				"Main Content",
			},
			excludes: []string{
				"<nav", "Navigation Menu",
				"<footer", "Site Footer",
				"<style", "color: red",
				"color: blue",
			},
		},
		{
			name: "keep styles only",
			opts: &CleaningOptions{
				KeepStyles: true,
			},
			contains: []string{
				"Main Content",
				`color: blue`,
				".header { color: red; }",
			},
			excludes: []string{
				"<header", "Site Header",
				"<nav", "Navigation Menu",
				"<footer", "Site Footer",
			},
		},
		{
			name: "keep all",
			opts: &CleaningOptions{
				KeepHeader:   true,
				KeepFooter:   true,
				KeepNav:      true,
				KeepStyles:   true,
				KeepComments: true,
			},
			contains: []string{
				"<header>Site Header</header>",
				"<nav>Navigation Menu</nav>",
				"<footer>Site Footer</footer>",
				`color: blue`,
				".header { color: red; }",
				"<!-- Comment -->",
			},
			excludes: []string{
				"<script", "document.write",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(CleanHTML([]byte(html), tt.opts))

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
