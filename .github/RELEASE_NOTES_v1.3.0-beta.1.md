# 🚀 gohome v1.3.0-beta.1 - Cobra/Viper Migration (Beta)

## ⚠️ Beta Release - Testing Required

This is a **beta release** for community testing. While extensively tested (45.1% coverage with 100% for critical packages), we recommend:
- ✅ Testing in non-production environments first
- 🐛 Reporting any issues you encounter
- 💬 Providing feedback on new features

---

## 🎯 What's New

### Modern CLI Framework
- **Cobra/Viper Integration:** Professional CLI architecture with subcommands
- **Better UX:** Organized help text with ASCII art logo
- **Improved Validation:** Better flag binding and error handling

### Multi-Format Configuration
- **JSON, YAML, TOML:** Use your preferred config format
- **New Location:** `~/.config/gohome/config.{json,yaml,toml}`
- **Backward Compatible:** Old `~/.gohome.json` still works

### Environment Variables
- **GOHOME_* Prefix:** Configure everything via env vars
- **Priority System:** CLI flags > env vars > config file > defaults
- **Examples:**
  ```bash
  export GOHOME_WORKSPACE="$HOME/projects"
  export GOHOME_PERIOD_DAYS=1
  export GOHOME_FORMAT="table"
  ```

### Config Management
New `gohome config` subcommand:
```bash
gohome config init   # Initialize with defaults
gohome config show   # Display current config
gohome config edit   # Edit in default editor
```

### Shell Completions
Auto-completion for all major shells:
```bash
gohome completion bash      # Bash
gohome completion zsh       # Zsh
gohome completion fish      # Fish
gohome completion powershell # PowerShell
```

### Comprehensive Testing
- 📊 **45.1% Overall Coverage** (up from 14.6%)
- ✅ **100% Coverage:** parser, git, renderer packages
- 🔒 **Security Tests:** Command injection prevention, input sanitization
- 🧪 **50+ Test Cases:** Table-driven tests with edge cases

---

## 📦 What to Test

### 1. Migration from v1.2
```bash
# Backup your config
cp ~/.gohome.json ~/.gohome.json.backup

# Install beta
go install github.com/anIcedAntFA/gohome/cmd/gohome@v1.3.0-beta.1

# Verify it works with old config
gohome report

# Try new config location (optional)
gohome config init
```

### 2. New Features
```bash
# Config management
gohome config show
gohome config edit

# Environment variables
export GOHOME_FORMAT="markdown"
gohome report

# Shell completions
gohome completion bash > /etc/bash_completion.d/gohome
# Restart shell and try: gohome <TAB>

# Different config formats
cp ~/.gohome.json ~/.config/gohome/config.yaml
# Edit to YAML format, test
```

### 3. Backward Compatibility
- ✅ Old config location (`~/.gohome.json`) still works
- ✅ All CLI flags unchanged
- ✅ Existing workflows continue without modification

### 4. Report Issues
Check for:
- Config loading problems
- Flag parsing errors
- Shell completion issues
- Environment variable conflicts
- Any unexpected behavior

---

## 📥 Installation

### Using go install (Recommended for Beta)
```bash
go install github.com/anIcedAntFA/gohome/cmd/gohome@v1.3.0-beta.1
```

### Download Binary
Download from [Release Assets](#assets) below:
- **Linux:** `gohome_1.3.0-beta.1_linux_x86_64.tar.gz`
- **macOS:** `gohome_1.3.0-beta.1_darwin_x86_64.tar.gz` (Intel) or `darwin_arm64` (M1/M2)
- **Windows:** `gohome_1.3.0-beta.1_windows_x86_64.zip`

```bash
# Example for Linux
wget https://github.com/anIcedAntFA/gohome/releases/download/v1.3.0-beta.1/gohome_1.3.0-beta.1_linux_x86_64.tar.gz
tar -xzf gohome_1.3.0-beta.1_linux_x86_64.tar.gz
sudo mv gohome /usr/local/bin/
```

---

## 📖 Documentation

- 📘 [Migration Guide](https://github.com/anIcedAntFA/gohome/blob/main/docs/v1.3_MIGRATION_GUIDE.md) - Complete v1.2 → v1.3 upgrade guide
- 📗 [CLI Guide](https://github.com/anIcedAntFA/gohome/blob/main/docs/CLI_GUIDE.md) - Updated command reference
- 📕 [CHANGELOG](https://github.com/anIcedAntFA/gohome/blob/main/CHANGELOG.md) - Full changelog

---

## 🐛 Known Issues

**None currently** - please report any you find!

---

## 💬 Feedback & Bug Reports

**Please report issues at:** https://github.com/anIcedAntFA/gohome/issues/new

Include:
- 🖥️ **OS and architecture** (Linux/macOS/Windows, x64/arm64)
- 🔧 **Go version** (if building from source)
- 📝 **Steps to reproduce**
- ✅ **Expected behavior**
- ❌ **Actual behavior**
- 📋 **Config file** (sanitized, without sensitive data)

---

## 🗓️ Timeline

- **Beta Testing:** 1-2 weeks (until ~Feb 1, 2026)
- **Stable Release:** v1.3.0 (after community feedback)

---

## 🙏 Thanks

Thank you for testing v1.3.0-beta.1! Your feedback helps make gohome better.

---

## 📊 Full Changelog

See [CHANGELOG.md](https://github.com/anIcedAntFA/gohome/blob/main/CHANGELOG.md#130-beta1---2026-01-18) for complete details.

**Full Diff:** https://github.com/anIcedAntFA/gohome/compare/v1.2.0...v1.3.0-beta.1
