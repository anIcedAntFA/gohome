# gohome

A fast, configurable Git standup & activity reporting CLI written in Go.

## Installation

```bash
npm install -g @ngockhoi96/gohome
```

Or use npx without installation:

```bash
npx @ngockhoi96/gohome --help
```

## Quick Start

```bash
# Generate report for today
gohome

# Show last 3 days
gohome --days 3

# Custom workspace path
gohome --path ~/projects

# Copy to clipboard
gohome --copy
```

## Features

- **🚀 Auto-Discovery:** Recursively finds git repositories in your workspace
- **⚡ Concurrency:** Scans multiple repos in parallel using Goroutines
- **🎨 Rich Output:** Supports multiple formats (text, table) and styles
- **📋 Clipboard Ready:** Copy reports directly to your clipboard
- **📝 Custom Tasks:** Add manual tasks alongside git commits
- **⚙️ Smart Config:** Persist preferences via `~/.gohome.json` or CLI flags

## Documentation

Full documentation available at: https://github.com/anIcedAntFA/gohome

## License

MIT License - see LICENSE file for details
