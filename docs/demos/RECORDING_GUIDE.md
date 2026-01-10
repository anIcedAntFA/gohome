# 🎬 Demo Recording Guide

This guide shows how to create terminal demo recordings for gohome using **VHS** (Video Hacking System) by Charm.

## Why VHS?

**VHS** allows you to write terminal GIFs **as code** instead of recording actual terminal sessions. This provides several advantages:

| Feature | Traditional Recording | VHS (Code-based) |
|---------|----------------------|------------------|
| **Reproducibility** | ❌ System-dependent | ✅ 100% deterministic |
| **Version Control** | ⚠️ Binary files only | ✅ Source code (`.tape`) |
| **Editing** | ⚠️ Limited (trim only) | ✅ Full control (every keystroke) |
| **CI/CD** | ⚠️ Manual process | ✅ GitHub Action built-in |
| **Testing** | ❌ | ✅ Golden files (`.ascii`) |
| **Consistency** | ⚠️ Varies by run | ✅ Pixel-perfect every time |

**Used by**: Charm's BubbleTea, Glow, Glamour, and hundreds of professional CLI tools.

## 📦 Prerequisites

Install VHS and its dependencies (Arch Linux):

```bash
# VHS requires ttyd and ffmpeg
yay -S vhs ttyd ffmpeg

# Verify installation
vhs --version  # Should show 0.10+
ttyd --version # Should show 1.7+
ffmpeg -version | head -1
```

**Other platforms:**
- macOS: `brew install vhs`
- Ubuntu/Debian: See [VHS installation docs](https://github.com/charmbracelet/vhs#installation)
- Windows: `scoop install vhs` or `winget install charmbracelet.vhs`

---

## 🎥 Creating Demos with VHS

### Quick Start

VHS uses `.tape` files (DSL) to script terminal sessions:

```bash
# Create a new tape file
vhs new demo.tape

# Edit the tape file
vim docs/demos/quickstart.tape

# Generate the GIF
vhs docs/demos/quickstart.tape
```

### Tape File Structure

A `.tape` file consists of commands that VHS executes:

```elixir
# Output configuration
Output docs/demos/quickstart.gif

# Terminal settings
Set Shell fish
Set FontSize 16
Set Width 1200
Set Height 600
Set Theme "Dracula"
Set TypingSpeed 50ms

# Script your demo
Type "gohome --version"
Enter
Sleep 2s

Type "gohome --help"
Enter
Sleep 3s
```

### Key VHS Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `Output` | Specify output file | `Output demo.gif` |
| `Set <Setting>` | Configure terminal | `Set FontSize 16` |
| `Type "<text>"` | Simulate typing | `Type "echo hello"` |
| `Enter` | Press Enter key | `Enter` |
| `Sleep <time>` | Wait/pause | `Sleep 2s` |
| `Ctrl+<key>` | Control sequences | `Ctrl+C` |
| `Hide`/`Show` | Toggle recording | `Hide` (for setup) |
| `Screenshot` | Capture frame | `Screenshot out.png` |

### Available Themes

```bash
# List all themes
vhs themes

# Popular themes:
- Dracula (purple/pink, used by gohome)
- Catppuccin Mocha
- Nord
- Tokyo Night
- Gruvbox Dark
- Monokai
```

### Terminal Settings Reference

```elixir
Set Shell fish              # Shell to use
Set FontSize 16             # Font size in pixels
Set FontFamily "JetBrains Mono"
Set Width 1200              # Terminal width in pixels
Set Height 600              # Terminal height in pixels
Set Padding 20              # Inner padding
Set Margin 10               # Outer margin
Set MarginFill "#282a36"    # Margin background color
Set BorderRadius 8          # Rounded corners
Set WindowBar Colorful      # Window controls style
Set TypingSpeed 50ms        # Typing delay (global)
Set Theme "Dracula"         # Color theme
Set FrameRate 60            # FPS for recording
Set PlaybackSpeed 1.0       # Speed multiplier
Set LoopOffset 0            # GIF loop start frame
Set CursorBlink true        # Blinking cursor
```

---

## 📝 Demo Examples

### Example 1: Quick Start Demo

File: `docs/demos/quickstart.tape`

```elixir
Output docs/demos/quickstart.gif

Set Shell fish
Set FontSize 16
Set Width 1200
Set Height 600
Set Theme "Dracula"
Set TypingSpeed 50ms

Type "# 🚀 gohome - Git Activity Aggregator"
Enter
Sleep 1s

Type "gohome --version"
Sleep 500ms
Enter
Sleep 2s

Type "gohome"
Enter
Sleep 3s
```

### Example 2: Configuration Demo

File: `docs/demos/config.tape`

```elixir
Output docs/demos/config.gif

Set Shell fish
Set Theme "Dracula"
Set TypingSpeed 50ms

Type `gohome -t "Code review: PR #123"`
Sleep 500ms
Enter
Sleep 3s

Type "gohome --format table --style markdown"
Enter
Sleep 3s

Type "gohome --icon --scope"
Enter
Sleep 3s
```

---

## 🚀 Workflow

### Using Makefile (Recommended)

```bash
# Generate all demos
make demo-record

# Validate tape files
make demo-validate

# List available themes
make demo-themes

# Clean generated files
make demo-clean
```

### Manual Execution

```bash
# Generate single demo
vhs docs/demos/quickstart.tape

# Watch mode (auto-regenerate on file change)
vhs docs/demos/quickstart.tape --watch

# Validate syntax
vhs validate docs/demos/quickstart.tape

# Output to multiple formats
# Edit .tape file:
Output demo.gif
Output demo.mp4
Output demo.webm
```

---

## 📏 Best Practices

### Writing Good Tape Files

1. **Use comments liberally:**
   ```elixir
   # Show version information
   Type "gohome --version"
   Enter
   Sleep 2s
   ```

2. **Set consistent terminal size:**
   ```elixir
   Set Width 1200
   Set Height 600
   ```

3. **Add natural pauses:**
   ```elixir
   Type "gohome"
   Sleep 500ms  # Brief pause before Enter
   Enter
   Sleep 3s     # Let output display
   ```

4. **Use Hide/Show for setup:**
   ```elixir
   Hide
   Type "make build && clear"  # Build binary (hidden)
   Enter
   Show
   Type "gohome --version"  # This appears in GIF
   ```

5. **Keep demos focused:**
   - One concept per demo
   - 30-60 seconds ideal
   - Show output clearly

### File Organization

```
docs/demos/
├── RECORDING_GUIDE.md          # This file
├── quickstart.tape             # Quick start demo script
├── quickstart.gif              # Generated GIF
├── config.tape                 # Configuration demo script
├── config.gif                  # Generated GIF
└── themes.tape                 # (optional) Theme showcase
```

### Version Control

VHS `.tape` files are **source code** - commit them to git:

```bash
# Commit tape files (always)
git add docs/demos/*.tape

# Commit GIFs (with Git LFS if >1MB)
git add docs/demos/*.gif
```

---

## 🧪 Integration Testing

VHS can generate `.ascii` files for golden testing:

```elixir
# In your .tape file
Output golden.ascii

# Run test
vhs test.tape

# Compare with committed golden file
diff golden.ascii golden.expected.ascii
```

Use in CI to ensure CLI output doesn't change unexpectedly.

---

## 🎯 CI/CD with GitHub Actions

Use the official VHS GitHub Action:

```yaml
name: Generate Demos

on: [push, pull_request]

jobs:
  demos:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: charmbracelet/vhs-action@v2
        with:
          path: 'docs/demos/*.tape'
```

This auto-generates GIFs on every commit!

---

## 🎨 Advanced Techniques

### Dynamic Content

```elixir
# Use environment variables
Env PROJECT_NAME "gohome"
Type "echo $PROJECT_NAME"
Enter
```

### Conditional Flow

```elixir
# Wait for output before continuing
Type "gohome --days 7"
Enter
Wait /Repository/ # Wait for "Repository" to appear
Sleep 1s
```

### Multiple Outputs

```elixir
# Generate multiple formats at once
Output demo.gif
Output demo.mp4
Output demo.webm
Output frames/  # PNG sequence
```

### Source Reuse

```elixir
# Reuse common settings
Source config.tape  # Import settings from another file
Type "gohome"
```

---

## 🔗 Resources

- [VHS GitHub](https://github.com/charmbracelet/vhs)
- [VHS Documentation](https://github.com/charmbracelet/vhs#vhs-command-reference)
- [VHS GitHub Action](https://github.com/charmbracelet/vhs-action)
- [Charm Community](https://charm.sh/)
- [Example Tapes](https://github.com/charmbracelet/vhs/tree/main/examples)

---

## 🆚 Comparison: VHS vs asciinema

| Aspect | asciinema + agg | VHS |
|--------|----------------|-----|
| **Workflow** | Record → Convert | Write code → Generate |
| **File format** | `.cast` (JSON lines) | `.tape` (DSL) |
| **Reproducibility** | ❌ System-dependent | ✅ Deterministic |
| **Editing** | Limited (trim only) | Full (modify any line) |
| **Version control** | Binary diffs | Text diffs |
| **CI/CD** | Manual setup | GitHub Action |
| **Testing** | No | Yes (golden files) |
| **Learning curve** | Low (just record) | Medium (learn DSL) |
| **Output quality** | Good | Excellent |
| **File size** | ~200KB | ~200KB (similar) |

**Verdict**: VHS is better for professional projects requiring reproducibility and CI/CD integration.

---

## 🚀 Quick Reference

### Generate Demos

```bash
# Using Makefile (recommended)
make demo-record        # Generate all GIFs
make demo-validate      # Check tape syntax
make demo-themes        # List available themes
make demo-clean         # Remove generated files

# Manual
vhs docs/demos/quickstart.tape
vhs docs/demos/config.tape --watch  # Auto-reload on change
```

### Common VHS Commands

```elixir
# Output
Output demo.gif

# Settings
Set Theme "Dracula"
Set FontSize 16
Set Width 1200
Set Height 600
Set TypingSpeed 50ms

# Actions
Type "command here"
Type@100ms "slow typing"
Enter
Sleep 2s
Ctrl+C

# Control
Hide  # Stop recording
Show  # Resume recording

# Wait for conditions
Wait /pattern/  # Wait for regex match
```

### Debugging Tape Files

```bash
# Validate syntax
vhs validate demo.tape

# Verbose output
VHS_DEBUG=1 vhs demo.tape

# Check single frame
vhs demo.tape --screenshot
```
