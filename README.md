# md-fetch

A powerful command-line tool that fetches web content and converts it to clean, readable Markdown format.

## Key Features

- **Bypass Anti-Scraping Measures**: Uses real browsers in headless mode to bypass 403 errors and CAPTCHAs that typically block programmatic scraping
- **Multiple Browser Support**: Uses Chrome, Firefox, or curl to fetch web content
- **Smart HTML Cleaning**: Removes unwanted JavaScript, CSS, and metadata while preserving content
- **JavaScript Support**: Properly renders JavaScript-heavy websites using Chrome or Firefox
- **Clean Markdown Output**: Converts cleaned HTML to well-formatted Markdown
- **AI/LLM Optimized**: Produces lightweight, clean text that's perfect for feeding into AI models

## Perfect for AI/LLM Applications

md-fetch is especially valuable for AI and Large Language Model (LLM) applications:
- **Clean Input**: Removes noise like scripts, styles, and metadata that could confuse LLMs
- **Token Efficiency**: Outputs lightweight Markdown, reducing token usage when feeding content to AI models
- **Context Preservation**: Maintains important content structure while eliminating irrelevant elements
- **Consistent Format**: Provides uniformly formatted text regardless of the source website's structure
- **Easy Integration**: Perfect for automating web content collection for AI training or real-time querying

## Why Use Real Browsers?

Many modern websites implement anti-scraping measures that block traditional HTTP requests:
- Return 403 Forbidden errors
- Present CAPTCHAs
- Require JavaScript execution
- Check for browser fingerprints

md-fetch solves this by using real browsers (Chrome/Firefox) in headless mode, which:
- Appears as a legitimate browser
- Executes JavaScript properly
- Handles modern web features
- Maintains your existing browser session

## Installation

### Pre-built Binaries
Visit our [releases page](https://github.com/nathabonfim59/md-fetch/releases) to download pre-built binaries:

#### Linux
- Standard build (dynamically linked)
- Musl build (statically linked, ideal for Alpine Linux and other musl-based systems)
- Debian/Ubuntu package (.deb)
- Red Hat/Fedora package (.rpm)

#### macOS
- Intel (amd64)
- Apple Silicon (arm64)

#### Windows
- 64-bit (amd64)

#### Package Installation
For Debian/Ubuntu:
```bash
sudo dpkg -i md-fetch_<version>_amd64.deb
```

For Red Hat/Fedora:
```bash
sudo rpm -i md-fetch-<version>.x86_64.rpm
```

### Using go install
```bash
go install github.com/nathabonfim59/md-fetch@latest
```

### From Source
1. Ensure you have Go 1.16 or later installed
2. Clone the repository:
   ```bash
   git clone https://github.com/nathabonfim59/md-fetch.git
   cd md-fetch
   ```
3. Build the project:
   ```bash
   go build -o bin/md-fetch main.go
   ```

## Usage

### CLI Mode

Basic usage:
```bash
md-fetch [url]
```

Options:
```bash
md-fetch [flags] [url]

Flags:
  -b, --browser string    Browser to use (optional, defaults to chrome > firefox > curl)
  -f, --filename string   Custom filename to save the content (optional, defaults to slugified URL)
  -h, --help             help for md-fetch
  -s, --save             Save content to a file with slugified URL name
```

Examples:
```bash
# Use default browser (tries Chrome, then Firefox, then curl)
md-fetch https://www.google.com

# Use Chrome specifically
md-fetch --browser chrome https://www.google.com

# Use Firefox and save to file
md-fetch --browser firefox --save https://www.google.com

# Use curl and specify custom filename
md-fetch --browser curl --save --filename output.md https://www.google.com
```

### Server Mode

Start the HTTP server:
```bash
md-fetch serve [flags]

Flags:
  -h, --help        help for serve
  -p, --port int    Port for HTTP server (default 8080)
```

The server provides a REST API for fetching content:

#### Single URL Request
```bash
curl -X POST http://localhost:8080/fetch \
  -H "Content-Type: application/json" \
  -d '{
    "urls": ["https://www.example.com"],
    "browser": "chrome"
  }'
```

#### Batch Request (Parallel Processing)
```bash
curl -X POST http://localhost:8080/fetch \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://www.example.com",
      "https://www.google.com",
      "https://www.github.com"
    ],
    "browser": "chrome"
  }'
```

#### OpenAPI Specification
Access the OpenAPI specification at:
```
http://localhost:8080/openapi.yaml
```

#### Response Format
```json
{
  "results": {
    "https://www.example.com": "# Example Domain\n\nThis domain is for...",
    "https://www.google.com": "# Google\n\n[Gmail](https://mail.google.com)..."
  },
  "errors": {
    "https://invalid.url": "error message"
  }
}
```

## Browser Support

The tool supports multiple browsers in the following priority order:

1. **Chrome/Chromium**: Best for JavaScript-heavy sites (default)
2. **Firefox**: Good alternative for JavaScript support
3. **curl**: Fallback for static content

You can specify which browser to use with the `-browser` flag.

## HTML Cleaning Features

- Removes JavaScript code:
  - Anonymous functions and IIFEs
  - Event listeners
  - Window assignments
  - Variable declarations
  - Google-specific scripts
  - MediaWiki RLQ functions

- Cleans CSS content:
  - Inline styles
  - Style blocks
  - Media queries
  - CSS definitions

- Removes metadata:
  - JSON-LD data
  - Schema.org markup
  - Configuration objects

## Project Structure

```
md-fetch/
├── cmd/                    # Command-line interface
├── internal/              
│   ├── browser/           # Browser implementations
│   │   ├── chrome.go      # Chrome/Chromium support
│   │   ├── firefox.go     # Firefox support
│   │   ├── curl.go        # curl support
│   │   └── html_cleaner.go # HTML cleaning logic
│   ├── converter/         # HTML to Markdown conversion
│   └── fetcher/           # Content fetching coordination
├── bin/                   # Compiled binaries
└── main.go                # Entry point
```

## Development

### Running Tests

```bash
go test ./...
```

### Adding New Features

1. **New Browser Support**: Implement the `Browser` interface in `internal/browser/browser.go`
2. **HTML Cleaning Rules**: Add patterns to `html_cleaner.go`
3. **Markdown Conversion**: Enhance `internal/converter/markdown.go`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
