# Server Mode

Start the HTTP server to process multiple URLs at once or integrate with other services:
```bash
md-fetch serve [flags]

Flags:
  -p, --port int    Port for HTTP server (default 8080)
```

## REST API Usage

The server provide a REST API for fetching content.

### Single or Batch URL Request

```bash
curl -X POST http://localhost:8080/fetch 
  -H "Content-Type: application/json" 
  -d '{
    "urls": ["https://www.example.com", "https://www.google.com"],
    "browser": "chrome"
  }'
```

### Response Format

```json
{
  "results": {
    "https://www.example.com": "# Example Domain

This domain is for...",
    "https://www.google.com": "# Google

[Gmail](https://mail.google.com)..."
  },
  "errors": {
    "https://invalid.url": "error message"
  }
}
```

## OpenAPI Specification

Access the interactive documentation or the JSON spec:
```bash
# Get the raw spec
curl http://localhost:8080/openapi.yaml
```

You can also view the full [API Explorer](https://nathabonfim59.github.io/md-fetch/site/) on the project website.
