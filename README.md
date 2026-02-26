# md-fetch

A powerful command-line tool that fetches web content and converts it to clean, readable Markdown format, optimized for AI and LLMs.

[![Docs](https://img.shields.io/badge/docs-site-blue)](https://nathabonfim59.github.io/md-fetch/site/)
[![License](https://img.shields.io/github/license/nathabonfim59/md-fetch)](LICENSE)

## Key Features

- **Bypass Anti-Scraping**: Uses real headless browsers (Chrome/Firefox) or curl to get content.
- **AI-Optimized**: Strips noise (scripts, styles, ads) to produce token-efficient Markdown.
- **Server Mode**: REST API for single or batch URL processing with parallel execution.
- **Agent Ready**: Ships with a specialized [Agent Skill](skills/md-fetch-cli/SKILL.md) for OpenCode, Gemini CLI, and Claude Code.

[Learn more about Features & Why Use Real Browsers →](docs/features.md)

## Quick Start

### Installation

```bash
# Using Go
go install github.com/nathabonfim59/md-fetch@latest

# As an Agent Skill
npx skills add nathabonfim59/md-fetch
```

[Detailed Installation Guide (Releases, Source, etc.) →](docs/installation.md)

### Basic Usage

```bash
# Fetch a URL to terminal
md-fetch https://example.com

# Save as Markdown file
md-fetch --save https://example.com

# Use a specific browser
md-fetch --browser firefox https://example.com
```

[Browser Support & Troubleshooting →](docs/browsers.md)

### Server Mode

```bash
# Start REST API server
md-fetch serve --port 8080
```

[Full Server Mode & API Guide →](docs/server.md)

## Documentation

- **[Full Landing Page & API Explorer](https://nathabonfim59.github.io/md-fetch/site/)**
- **[CLI Reference](skills/md-fetch-cli/references/COMMANDS.md)**
- **[Agent Skill Guide](skills/md-fetch-cli/SKILL.md)**
- **[Development & Contributing](docs/development.md)**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
