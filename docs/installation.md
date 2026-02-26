# Installation

## Pre-built Binaries

Visit our [releases page](https://github.com/nathabonfim59/md-fetch/releases) to download pre-built binaries for your platform.

### Linux
- **Standard**: Dynamically linked.
- **Musl**: Statically linked (ideal for Alpine Linux).
- **Packages**: .deb (Debian/Ubuntu) and .rpm (Red Hat/Fedora).

### macOS
- **Intel (amd64)**
- **Apple Silicon (arm64)**

### Windows
- **64-bit (amd64)**

### Package Installation Example

**Debian/Ubuntu:**
```bash
sudo dpkg -i md-fetch_<version>_amd64.deb
```

**Red Hat/Fedora:**
```bash
sudo rpm -i md-fetch-<version>.x86_64.rpm
```

## Using Go

If you have Go installed (1.16+):
```bash
go install github.com/nathabonfim59/md-fetch@latest
```

## Agent Skill (OpenCode, Gemini CLI, Claude Code)

Install the specialized agent skill to use md-fetch with compatible AI agents:
```bash
npx skills add nathabonfim59/md-fetch
```

## From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/nathabonfim59/md-fetch.git
   cd md-fetch
   ```
2. Build the project:
   ```bash
   go build -o bin/md-fetch main.go
   ```
