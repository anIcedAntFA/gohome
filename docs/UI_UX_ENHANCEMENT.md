# 🎨 UI/UX Enhancement with Lip Gloss

**Status:** Planning Phase  
**Target Release:** v1.4.0 (Phase 2)  
**Priority:** High - Foundation for visual improvements

## 📋 Overview

Transform gohome's terminal output from functional to beautiful using [Lip Gloss](https://github.com/charmbracelet/lipgloss), the powerful styling library from Charm. This enhancement will make daily standup reports more visually appealing, easier to scan, and professional-looking while maintaining readability and accessibility.

## 🎯 Goals

1. **Visual Hierarchy:** Clear distinction between sections (repo name, commits, tasks)
2. **Color Consistency:** Cohesive color palette that works in light/dark terminals
3. **Branding:** Distinctive gohome identity through styled headers and banners
4. **Accessibility:** Maintain readability with proper contrast and adaptive colors
5. **Performance:** Zero impact on scan/render speed
6. **Foundation:** Prepare for future interactive TUI (Phase 3)

## 🎨 Design System

### Color Palette

Following terminal-friendly principles with adaptive light/dark support:

```go
// Brand Colors
var (
    // Primary brand color - used for headers, important elements
    PrimaryColor = lipgloss.AdaptiveColor{
        Light: "#5A67D8", // Indigo-600
        Dark:  "#818CF8", // Indigo-400
    }
    
    // Accent color - used for highlights, success states
    AccentColor = lipgloss.AdaptiveColor{
        Light: "#059669", // Emerald-600
        Dark:  "#34D399", // Emerald-400
    }
    
    // Muted color - used for secondary text, metadata
    MutedColor = lipgloss.AdaptiveColor{
        Light: "#6B7280", // Gray-500
        Dark:  "#9CA3AF", // Gray-400
    }
    
    // Danger color - used for errors, warnings
    DangerColor = lipgloss.AdaptiveColor{
        Light: "#DC2626", // Red-600
        Dark:  "#F87171", // Red-400
    }
)

// Commit Type Colors (Conventional Commits)
var CommitTypeColors = map[string]lipgloss.Color{
    "feat":     "#10B981", // Green
    "fix":      "#EF4444", // Red
    "docs":     "#3B82F6", // Blue
    "style":    "#8B5CF6", // Purple
    "refactor": "#F59E0B", // Amber
    "perf":     "#EC4899", // Pink
    "test":     "#06B6D4", // Cyan
    "chore":    "#6B7280", // Gray
    "ci":       "#14B8A6", // Teal
    "build":    "#F97316", // Orange
}
```

### Typography Styles

```go
// Banner/Logo - Displayed with `gohome -h` or `gohome version`
var BannerStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(PrimaryColor).
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(AccentColor).
    Padding(1, 2).
    MarginBottom(1)

// Section Headers (e.g., "📂 Repositories", "📝 Tasks")
var HeaderStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(PrimaryColor).
    BorderStyle(lipgloss.ThickBorder()).
    BorderBottom(true).
    BorderForeground(AccentColor).
    MarginTop(1).
    MarginBottom(1)

// Repository Names
var RepoNameStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(AccentColor).
    PaddingLeft(2)

// Commit Messages
var CommitMessageStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("252")). // Light gray
    PaddingLeft(4)

// Commit Types (feat, fix, etc.)
var CommitTypeStyle = lipgloss.NewStyle().
    Bold(true).
    Padding(0, 1).
    MarginRight(1)

// Metadata (branch names, timestamps, author)
var MetadataStyle = lipgloss.NewStyle().
    Foreground(MutedColor).
    Italic(true)

// Task Items
var TaskStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#FBBF24")). // Amber
    PaddingLeft(2)

// Empty State
var EmptyStateStyle = lipgloss.NewStyle().
    Foreground(MutedColor).
    Italic(true).
    Align(lipgloss.Center)
```

## 🖼️ Enhanced Output Examples

### Current Output (Plain)
```
🔍 Scanning repositories...
✅ Found 3 repositories

📂 gohome (main)
  feat: add whitelist feature
  fix: resolve config merge bug
  docs: update README

📂 my-app (develop)
  feat: implement user authentication
  test: add login tests

📝 Custom Tasks:
  - Attended team standup
  - Reviewed 2 PRs
```

### Enhanced Output (with Lip Gloss)
```
╭─────────────────────────────────────────────╮
│                                             │
│   ██████╗  ██████╗ ██╗  ██╗ ██████╗ ███╗   ███╗███████╗  │
│  ██╔════╝ ██╔═══██╗██║  ██║██╔═══██╗████╗ ████║██╔════╝  │
│  ██║  ███╗██║   ██║███████║██║   ██║██╔████╔██║█████╗    │
│  ██║   ██║██║   ██║██╔══██║██║   ██║██║╚██╔╝██║██╔══╝    │
│  ╚██████╔╝╚██████╔╝██║  ██║╚██████╔╝██║ ╚═╝ ██║███████╗  │
│   ╚═════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝  │
│                                             │
│          Daily Standup Report               │
│           2026-01-20 • @ngockhoi96          │
╰─────────────────────────────────────────────╯

🔍 Scanning repositories...
✅ Found 3 repositories (scanned in 124ms)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📂 REPOSITORIES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

┌─ gohome ─────────────────────────────────┐
│ 📍 main • 3 commits • last 24h          │
└──────────────────────────────────────────┘
    ✨ feat: add whitelist feature
    🐛 fix: resolve config merge bug
    📚 docs: update README

┌─ my-app ─────────────────────────────────┐
│ 📍 develop • 2 commits • last 24h       │
└──────────────────────────────────────────┘
    ✨ feat: implement user authentication
    🧪 test: add login tests

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📝 CUSTOM TASKS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  • Attended team standup
  • Reviewed 2 PRs

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💡 Tip: Use --copy to copy this report to clipboard
```

### Table Output (Enhanced)
```
╭─────────────────────────────────────────────────────────────────╮
│                      DAILY STANDUP REPORT                       │
│                  2026-01-20 • @ngockhoi96                       │
╰─────────────────────────────────────────────────────────────────╯

┏━━━━━━━━━━━━━━━┳━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ REPOSITORY    ┃ TYPE  ┃ MESSAGE                              ┃
┡━━━━━━━━━━━━━━━╇━━━━━━━╇━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┩
│ gohome        │ ✨feat │ add whitelist feature                │
│               │ 🐛fix  │ resolve config merge bug             │
│               │ 📚docs │ update README                        │
├───────────────┼───────┼──────────────────────────────────────┤
│ my-app        │ ✨feat │ implement user authentication        │
│               │ 🧪test │ add login tests                      │
└───────────────┴───────┴──────────────────────────────────────┘

TASKS:
  • Attended team standup
  • Reviewed 2 PRs
```

## 🏗️ Implementation Plan

### Phase 1: Foundation (Week 1-2)

1. **Dependency Management**
   ```bash
   go get github.com/charmbracelet/lipgloss@latest
   ```

2. **Create Style Package**
   - File: `internal/ui/styles.go`
   - Define all style constants and color palette
   - Adaptive color support for light/dark terminals
   - Export style builder functions

3. **Renderer Refactoring**
   - Update `internal/renderer/printer.go`
   - Inject lipgloss styles into existing render logic
   - Maintain backward compatibility with `--format` flag
   - Add new `--style` flag: `classic`, `modern`, `minimal`

### Phase 2: Visual Components (Week 3-4)

4. **Banner Component**
   - ASCII art logo generator
   - Animated banner option (optional flag)
   - Version and metadata display

5. **Repository Card Component**
   - Boxed layout with borders
   - Branch indicator with color coding
   - Commit count badge
   - Time range indicator

6. **Commit List Component**
   - Icon/emoji rendering
   - Type-based color coding
   - Scope highlighting
   - Truncation with ellipsis for long messages

7. **Table Enhancement**
   - Replace `tablewriter` with `lipgloss/table`
   - Custom cell styling per row/column
   - Header styling
   - Alternating row colors

### Phase 3: Advanced Features (Week 5-6)

8. **Theme System**
   - Predefined themes: `default`, `ocean`, `forest`, `sunset`, `monochrome`
   - Custom theme from config file
   - `gohome config set theme ocean`

9. **Layout System**
   - Grid layout for multi-column output
   - Responsive width detection
   - Horizontal/vertical alignment utilities

10. **Progress Indicators**
    - Styled spinner during scan
    - Progress bars for long operations
    - Success/error state indicators

### Phase 4: Testing & Documentation (Week 7-8)

11. **Unit Tests**
    - Test style rendering with different terminal widths
    - Color degradation tests (TrueColor → ANSI256 → ANSI16)
    - Snapshot testing for output consistency

12. **Documentation**
    - Update CLI_GUIDE.md with style examples
    - Add screenshots/GIFs to README
    - Document theme customization
    - Add troubleshooting section for color issues

## 📦 Project Structure

```
internal/
  ui/
    styles.go         # Core style definitions and color palette
    banner.go         # ASCII banner renderer
    card.go           # Repository card component
    table.go          # Enhanced table renderer
    theme.go          # Theme management
    utils.go          # Width detection, alignment helpers
  renderer/
    printer.go        # Main renderer (updated to use ui package)
    printer_test.go   # Renderer tests
```

## 🎛️ Configuration

### Config File Addition

```json
{
  "ui": {
    "style": "modern",
    "theme": "default",
    "noBanner": false,
    "noColor": false,
    "emojiSet": "default"
  }
}
```

### New CLI Flags

```bash
--style <name>       # Output style: classic, modern, minimal (default: modern)
--theme <name>       # Color theme: default, ocean, forest, sunset, mono
--no-banner          # Disable ASCII banner
--no-color           # Disable all color output (for piping)
--emoji-set <name>   # Emoji set: default, nerd, ascii, none
```

### Environment Variables

```bash
GOHOME_STYLE=modern
GOHOME_THEME=ocean
GOHOME_NO_COLOR=1  # Standard convention
NO_COLOR=1          # Respect system-wide no-color
```

## 🚀 Future Enhancements (Phase 3)

1. **Interactive Mode Foundation**
   - Component-based architecture prepares for Bubble Tea integration
   - Reusable styled components in TUI

2. **Animated Elements**
   - Smooth transitions when switching views
   - Loading animations with progress
   - Confetti/celebration effects for milestones

3. **Custom ASCII Art**
   - User-provided banner from config
   - Project-specific logos per repository

4. **Rich Markdown Rendering**
   - Syntax highlighting for code blocks in commit messages
   - Link detection and styling
   - Integration with glamour library

## 📊 Success Metrics

- **Visual Appeal:** Positive user feedback on new UI
- **Accessibility:** Works correctly in light/dark terminals
- **Performance:** No measurable slowdown in render time
- **Adoption:** Users prefer new styles over classic
- **Compatibility:** Zero breaking changes to existing output formats

## 🔗 References

- [Lip Gloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Lip Gloss Examples](https://github.com/charmbracelet/lipgloss/tree/master/examples)
- [Charm Color Design System](https://github.com/charmbracelet)
- [Terminal Color Standards](https://en.wikipedia.org/wiki/ANSI_escape_code#Colors)

## 📝 Notes

- All color values should support adaptive light/dark modes
- Respect `NO_COLOR` environment variable convention
- Test on major terminals: iTerm2, Terminal.app, Windows Terminal, Alacritty, Kitty
- Consider accessibility guidelines (WCAG contrast ratios)
- Keep ASCII art optional for minimal environments
- Maintain JSON/Markdown export without styling
