---
name: md-fetch-cli
description: Use the md-fetch CLI to retrieve web content as clean Markdown, save output to files, or run the local HTTP API server. Activate when users ask to fetch web pages, convert pages to Markdown, batch fetch URLs via API, or troubleshoot browser selection and fetch failures in this repository.
compatibility: Requires the md-fetch binary and at least one supported fetch backend (chrome/chromium, firefox, or curl).
metadata:
  author: nathabonfim59
  version: "1.0"
---

# md-fetch CLI skill

Use this skill for tasks that operate the `md-fetch` command-line tool (not for general web browsing through unrelated tools).

## What this tool does

- Fetches `http/https` URLs and returns Markdown or plain text/JSON output.
- Uses browser backends in this priority when `--browser` is not set: `chrome`, `firefox`, `curl`.
- Can save output to a `.md` file.
- Can run as an HTTP service for single or batch URL fetches.

## Operating procedure

1. Confirm the request type:
   - Single fetch to terminal
   - Fetch and save to file
   - Start/operate API server mode
2. Always use the installed CLI executable:
   - `md-fetch <args>`
3. Prefer explicit browser selection when troubleshooting:
   - `--browser chrome`
   - `--browser firefox`
   - `--browser curl`
4. For file output, use `--save` and optionally `--filename <name>.md`.
5. For API mode, run `serve` and call `POST /fetch` with JSON body containing `urls` and optional `browser`.

## Canonical commands

```bash
# Single URL fetch
md-fetch https://example.com

# Force a browser backend
md-fetch --browser chrome https://example.com

# Save output with generated filename
md-fetch --save https://example.com

# Save output with explicit filename
md-fetch --save --filename example.md https://example.com

# Run server mode on default port 8080
md-fetch serve

# Run server mode on custom port
md-fetch serve --port 9090
```

## API usage

```bash
curl -X POST http://localhost:8080/fetch \
  -H "Content-Type: application/json" \
  -d '{
    "urls": ["https://example.com", "https://github.com"],
    "browser": "chrome"
  }'
```

OpenAPI document:

```text
http://localhost:8080/openapi.yaml
```

## Behavior notes

- If URL has no scheme, `https://` is automatically added.
- Supported explicit backends: `chrome` (or `chromium`), `firefox`, `curl`.
- JSON responses are pretty-printed and wrapped in fenced Markdown.
- Invalid method on `/fetch` returns `405`; invalid JSON body returns `400`.

## Troubleshooting checklist

1. "failed to initialize browser": install/verify a backend executable (`chrome/chromium`, `firefox`, or `curl`) on PATH.
2. "site cannot be reached": verify URL/network and retry with a different backend.
3. Empty/poor output: retry with `--browser chrome` for JS-heavy pages.
4. Save issues: ensure write permissions for target directory.

## References

See `references/COMMANDS.md` for concise command snippets.
