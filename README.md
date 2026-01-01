# gohome 🏠

> A blazing fast, cross-platform Git aggregation and reporting CLI tool.

![Go Version](https://img.shields.io/github/go-mod/go-version/anIcedAntFA/go-home)
![License](https://img.shields.io/github/license/anIcedAntFA/go-home)
[![Go Report Card](https://goreportcard.com/badge/github.com/anIcedAntFA/go-home)](https://goreportcard.com/report/github.com/anIcedAntFA/go-home)

**Forgot what you worked on yesterday?**

**gohome** automates your daily status reporting by recursively scanning your workspace to find git repositories. It aggregates commit logs from multiple projects instantly and formats them into beautiful, ready-to-share reports.

Perfect for **Daily Standups**, **Weekly Summaries**, or tracking your **Personal Coding Habits**.

## ✨ Features

- **🚀 Auto-Discovery:** Recursively finds git repositories in your workspace.
- **⚡ Concurrency:** Scans multiple repos in parallel using Goroutines for maximum speed.
- **🎨 Rich Output:** Supports Markdown tables, plain text, and custom styles (Nature, Tech, etc.).
- **📋 Clipboard Ready:** Copy reports directly to your system clipboard with `--copy`.
- **⚙️ Smart Config:** Persist your preferences via `~/.gohome.json` or use command-line aliases.

## 📦 Installation

### 🛠️ From Source (Go Developers)

If you have Go installed (1.20+), you can install the latest version directly:

```bash
go install github.com/anIcedAntFA/go-home/cmd/gohome@latest
```

Make sure your `$GOPATH/bin` is in your `$PATH`.

## 🚀 Usage

Simply run the tool in your workspace directory:

```bash
gohome
```

### 🧪 Common Examples

**1️⃣ Basic Usage (Last 1 day)**

```bash
gohome
```

**2️⃣ Look back 3 days**

```bash
gohome -d 3
```

**3️⃣ Generate a Table Report**

```bash
gohome -f table -s nature
```

**4️⃣ Copy to Clipboard**

This is useful for pasting directly into Slack/Teams/Discord:

```bash
gohome -d 1 --copy
```

**5️⃣ Save Settings**

Save your favorite flags as default (so you don't have to type them next time):

```bash
gohome -p /Users/ngockhoi96/workspace -d 1 -f table --save
```

## 🔧 Configuration

**gohome** looks for a config file at `~/.gohome.json`. You can create it manually or use the `--save` flag to auto-generate it.

### 🧾 Flags Reference

| Flag       | Alias | Description                                  | Default     |
| ---------- | ----- | -------------------------------------------- | ----------- |
| `--days`   | `-d`  | Number of days to look back                  | 1           |
| `--weeks`  | `-w`  | Number of weeks to look back                 | 0           |
| `--month`  | `-m`  | Number of months to look back                | 0           |
| `--years`  | `-y`  | Number of years to look back                 | 0           |
| `--path`   | `-p`  | Root path to scan                            | `.`         |
| `--author` | `-a`  | Git author name (auto-detected)              | System User |
| `--format` | `-f`  | Output format (text, table)                  | text        |
| `--style`  | `-s`  | Table style (normal, markdown, nature, tech) | normal      |
| `--copy`   | `-cp` | Copy output to clipboard                     | false       |
| `--icon`   | `-i`  | Show icon column (table only)                | false       |
| `--scope`  | `-c`  | Show scope column (table only)               | false       |
| `--save`   |       | Save current flags as default                | false       |
| `--help`   | `-h`  | Show help message                            |             |

## 🗺️ Roadmap

See [ROADMAP.md](ROADMAP.md) for the full development plan and upcoming features.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

### 🧑‍💻 Development Setup

```bash
# Clone the repo
git clone https://github.com/anIcedAntFA/go-home.git
cd go-home

# Install dependencies
go mod tidy

# Run locally
go run cmd/gohome/main.go
```

## ❤️ Credits & Motivation

**gohome** is heavily inspired by the awesome [git-standup](https://github.com/kamranahmedse/git-standup) utility by [Kamran Ahmed](https://github.com/kamranahmedse).

While `git-standup` is great, **gohome** was built to address specific personal needs for daily reporting, such as:

- **Rich formatting:** Tables, icons, and custom styles.
- **Workflow integration:** Direct clipboard support.
- **Smart config:** Persisted settings for zero-setup runs.

This project also serves as a practical journey to master **Go (Golang)**, implementing concepts like Concurrency, CLI architecture, and Cross-platform distribution.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
