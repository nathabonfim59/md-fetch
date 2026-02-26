# Features

## Key Features

- **Bypass Anti-Scraping Measures**: Uses real browsers in headless mode to bypass 403 errors and CAPTCHAs that typically block programmatic scraping.
- **Multiple Browser Support**: Uses Chrome, Firefox, or curl to fetch web content.
- **Smart HTML Cleaning**: Removes unwanted JavaScript, CSS, and metadata while preserving content.
- **JavaScript Support**: Properly renders JavaScript-heavy websites using Chrome or Firefox.
- **Clean Markdown Output**: Converts cleaned HTML to well-formatted Markdown.
- **AI/LLM Optimized**: Produces lightweight, clean text that's perfect for feeding into AI models.

## Perfect for AI/LLM Applications

md-fetch is especially valuable for AI and Large Language Model (LLM) applications:
- **Clean Input**: Removes noise like scripts, styles, and metadata that could confuse LLMs.
- **Token Efficiency**: Outputs lightweight Markdown, reducing token usage when feeding content to AI models.
- **Context Preservation**: Maintains important content structure while eliminating irrelevant elements.
- **Consistent Format**: Provides uniformly formatted text regardless of the source website's structure.
- **Easy Integration**: Perfect for automating web content collection for AI training or real-time querying.

## Why Use Real Browsers?

Many modern websites implement anti-scraping measures that block traditional HTTP requests:
- Return 403 Forbidden errors
- Present CAPTCHAs
- Require JavaScript execution
- Check for browser fingerprints

md-fetch solves this by using real browsers (Chrome/Firefox) in headless mode, which:
- Appears as a legitimate browser.
- Executes JavaScript properly.
- Handles modern web features.
- Maintains your existing browser session.

## HTML Cleaning Details

md-fetch strips:
- **JavaScript code**: Anonymous functions, IIFEs, event listeners, window assignments, etc.
- **CSS content**: Inline styles, style blocks, media queries.
- **Metadata**: JSON-LD, Schema.org markup, configuration objects.
