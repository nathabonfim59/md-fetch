# Development

## Project Structure

```
md-fetch
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

## Contributing

1. Fork the repository.
2. Create your feature branch: `git checkout -b feature/amazing-feature`.
3. Commit your changes: `git commit -m 'Add some amazing feature'`.
4. Push to the branch: `git push origin feature/amazing-feature`.
5. Open a Pull Request.

## Testing

Run tests with Go:
```bash
go test ./...
```

## Adding New Features

- **New Browser Support**: Implement the `Browser` interface in `internal/browser/browser.go`.
- **HTML Cleaning Rules**: Add patterns to `html_cleaner.go`.
- **Markdown Conversion**: Enhance `internal/converter/markdown.go`.
