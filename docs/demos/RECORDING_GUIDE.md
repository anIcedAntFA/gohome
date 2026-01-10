# 🎬 Demo Recording Guide

This guide shows how to create terminal demo recordings for gohome using asciinema.

## 📦 Prerequisites

Install required tools (Arch Linux):

```bash
# Terminal recorder
yay -S asciinema

# GIF generator (for GitHub embedding)
yay -S agg

# Alternative: SVG generator (smaller file size, scalable)
npm install -g svg-term-cli
```

---

## 🎥 Recording Sessions

### Demo 1: Quick Start (Installation + Basic Usage)

**Theme:** Show new users how to get started

**Script:**
```bash
# Start recording
asciinema rec docs/demos/quickstart.cast

# === Recording starts ===
# 1. Show version (pretend fresh install)
gohome --version

# 2. Show help
gohome --help

# 3. Run basic command (last 24 hours)
gohome

# 4. Run with custom time range
gohome --days 7

# Exit with Ctrl+D
```

### Demo 2: Configuration & Customization

**Theme:** Advanced features (config file, custom tasks, formats)

**Script:**
```bash
# Start recording
asciinema rec docs/demos/config.cast

# === Recording starts ===
# 1. Show current config
gohome --help | grep -A 20 "Configuration"

# 2. Save config to file
gohome --workspaces ~/projects --days 1 --format table --save

# 3. Add custom tasks
gohome -t "Code review: PR #123" -t "Team standup meeting"

# 4. Different output styles
gohome --style markdown
gohome --style nature

# 5. Copy to clipboard
gohome --copy
echo "Report copied to clipboard!"

# Exit with Ctrl+D
```

### Demo 3: Installation Process

**Theme:** Show install script in action

**Script:**
```bash
# Start recording
asciinema rec docs/demos/install.cast

# === Recording starts ===
# 1. Download and run install script
curl -fsSL https://raw.githubusercontent.com/anIcedAntFA/gohome/main/scripts/install.sh | bash

# 2. Verify installation
gohome --version
which gohome

# Exit with Ctrl+D
```

---

## 🎨 Converting to GIF/SVG

### Option A: GIF (Best for GitHub README)

```bash
# Convert .cast to .gif with agg
agg docs/demos/quickstart.cast docs/demos/quickstart.gif

# With custom settings
agg --theme monokai \
    --font-size 16 \
    --cols 100 \
    --rows 30 \
    --speed 1.5 \
    docs/demos/quickstart.cast \
    docs/demos/quickstart.gif
```

**agg options:**
- `--theme`: Color theme (monokai, dracula, nord, etc.)
- `--font-size`: Font size (default: 14)
- `--cols/--rows`: Terminal dimensions
- `--speed`: Playback speed multiplier (1.5 = 50% faster)
- `--idle-time-limit`: Max pause duration (e.g., 2.0 = max 2 seconds)

### Option B: SVG (Smaller size, scalable)

```bash
# Convert .cast to .svg
cat docs/demos/quickstart.cast | svg-term --out docs/demos/quickstart.svg

# With custom settings
cat docs/demos/quickstart.cast | svg-term \
    --width 100 \
    --height 30 \
    --at 9999 \
    --window \
    --out docs/demos/quickstart.svg
```

**svg-term options:**
- `--width/--height`: Terminal dimensions (columns/rows)
- `--at`: Timestamp to capture (9999 = end of recording)
- `--window`: Add window frame decoration
- `--from/--to`: Trim recording

---

## ✂️ Editing Recordings

### Trim recording (remove beginning/end)

```bash
# Remove first 2 seconds and last 1 second
asciinema cut --start 2 --end -1 docs/demos/raw.cast docs/demos/trimmed.cast
```

### Speed up playback

Edit the `.cast` file directly - each line has a timestamp:
```json
[0.123, "o", "text"]  // [timestamp, type, content]
```

Or use agg's `--speed` flag during conversion.

---

## 📏 Best Practices

### Recording Tips

1. **Set consistent terminal size:**
   ```bash
   # Before recording, resize terminal to 100x30
   printf '\e[8;30;100t'
   ```

2. **Use deliberate typing:**
   - Type at moderate speed (not too fast)
   - Add natural pauses between commands
   - Avoid typos (or edit them out)

3. **Keep it short:**
   - Aim for 30-60 seconds per demo
   - One concept per recording
   - Use `--idle-time-limit` to cut long pauses

4. **Add context:**
   - Show command output clearly
   - Include success/error messages
   - Demonstrate real-world use cases

### File Organization

```
docs/demos/
├── RECORDING_GUIDE.md          # This file
├── quickstart.cast             # Raw recording
├── quickstart.gif              # Converted GIF
├── config.cast
├── config.gif
├── install.cast
└── install.gif
```

### Git LFS for Large Files

If GIFs are >5MB, use Git LFS:

```bash
# Install Git LFS
yay -S git-lfs
git lfs install

# Track GIF files
git lfs track "docs/demos/*.gif"
git add .gitattributes
```

---

## 🚀 Publishing

### In README.md

```markdown
## 🎬 Quick Demo

![gohome demo](docs/demos/quickstart.gif)

Or link to asciinema.org:
[![asciicast](https://asciinema.org/a/RECORDING_ID.svg)](https://asciinema.org/a/RECORDING_ID)
```

### Upload to asciinema.org (optional)

```bash
# Upload recording to asciinema.org for embedding
asciinema upload docs/demos/quickstart.cast

# Get shareable link and embed code
```

---

## 🎯 Example Recording Session

```bash
# 1. Prepare terminal
clear
printf '\e[8;30;100t'  # Resize to 100x30

# 2. Start recording
asciinema rec docs/demos/quickstart.cast

# 3. Perform demo (type naturally, add pauses)
gohome --version
sleep 1
gohome --help
sleep 2
gohome
sleep 1

# 4. Stop recording (Ctrl+D)

# 5. Convert to GIF
agg --theme monokai \
    --font-size 16 \
    --speed 1.3 \
    --idle-time-limit 2 \
    docs/demos/quickstart.cast \
    docs/demos/quickstart.gif

# 6. Preview
xdg-open docs/demos/quickstart.gif

# 7. Commit and push
git add docs/demos/quickstart.{cast,gif}
git commit -m "📹 demo: add quickstart recording"
```

---

## 🔗 Resources

- [asciinema docs](https://docs.asciinema.org/)
- [agg GitHub](https://github.com/asciinema/agg)
- [svg-term-cli](https://github.com/marionebl/svg-term-cli)
- [Charm BubbleTea demos](https://github.com/charmbracelet/bubbletea#examples) (for inspiration)
