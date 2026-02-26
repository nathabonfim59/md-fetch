# Browser Support

The tool supports multiple browsers in the following priority order:

1. **Chrome/Chromium**: Best for JavaScript-heavy sites (default).
2. **Firefox**: Good alternative for JavaScript support.
3. **curl**: Fallback for static content.

You can specify which browser to use with the `--browser` (or `-b`) flag.

## Requirements

Ensure that the browser binary is installed and available in your system's `PATH`.

- **Chrome/Chromium**: `google-chrome`, `chromium-browser`, `chrome`, or `chromium`.
- **Firefox**: `firefox`.
- **curl**: `curl`.

## Troubleshooting

- **"failed to initialize browser"**: Check if the browser is installed and on your `PATH`.
- **Site cannot be reached**: Verify the URL and network connection.
- **Empty/poor output**: Try using `--browser chrome` for JS-heavy sites.
